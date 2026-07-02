package ecommerce.filter;

import io.opentelemetry.api.trace.Span;
import jakarta.servlet.Filter;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.ServletResponse;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.stereotype.Component;

import java.io.IOException;

@Component
public class TraceIdFilter implements Filter {

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
            throws IOException, ServletException {
        
        HttpServletResponse httpServletResponse = (HttpServletResponse) response;
        Span currentSpan = Span.current();
        
        if (currentSpan != null && currentSpan.getSpanContext().isValid()) {
            httpServletResponse.setHeader("X-Request-Id", currentSpan.getSpanContext().getTraceId());
        }
        
        chain.doFilter(request, response);
    }
}
