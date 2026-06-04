require 'bunny'
require 'json'

$stdout.sync = true

puts "Starting Notification Worker..."

rabbitmq_url = ENV['RABBITMQ_URL'] || 'amqp://guest:guest@localhost:5672'

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
    # Process notification logic here
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