require 'bunny'
require 'json'
require 'pg'

$stdout.sync = true

puts "Starting Notification Worker..."

rabbitmq_url = ENV['RABBITMQ_URL'] || 'amqp://guest:guest@localhost:5672'
database_url = ENV['DATABASE_URL'] || 'postgres://user:password@localhost:5432/notification_db'

def connect_db(url)
  puts "Connecting to DB..."
  conn = PG.connect(url)
  puts "Connected to DB"
  conn
rescue StandardError => e
  puts "DB Connection failed: #{e.message}. Retrying..."
  sleep 5
  retry
end

db_conn = connect_db(database_url)

def connect_with_retry(url)
  connection = Bunny.new(url)
  begin
    connection.start
    puts "Connected to RabbitMQ"
    return connection
  rescue StandardError => e
    puts "Connection failed: #{e.message}. Retrying..."
    sleep 5
    retry
  end
end

connection = connect_with_retry(rabbitmq_url)
sleep 2 # Small delay after connection

begin
  channel = connection.create_channel
  queue = channel.queue('notifications', durable: true)

  puts " [*] Waiting for messages in #{queue.name}. To exit press CTRL+C"
  queue.subscribe(block: true) do |delivery_info, properties, body|
    puts " [x] Received notification: #{body}"
    
    # Persist to DB
    begin
      db_conn.exec_params("INSERT INTO notifications (body) VALUES ($1)", [body])
      puts " [x] Notification persisted to DB"
    rescue StandardError => e
      puts "Error persisting to DB: #{e.message}"
      # Try to reconnect DB if needed
      db_conn = connect_db(database_url)
    end
  end
rescue Interrupt => _
  puts "Shutting down..."
  connection.close
rescue StandardError => e
  puts "Error: #{e.message}. Reconnecting..."
  sleep 2
  connection = connect_with_retry(rabbitmq_url)
  retry
end