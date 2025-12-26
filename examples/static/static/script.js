// Static File Server Example JavaScript

// Simple demonstration script for the static file server
document.addEventListener('DOMContentLoaded', function() {
    console.log('Static file server example loaded!');
    
    // Add some interactivity to the page
    const body = document.body;
    const colors = ['#667eea', '#764ba2', '#f093fb', '#f5576c', '#4facfe'];
    let colorIndex = 0;
    
    // Change background color every 5 seconds
    setInterval(function() {
        colorIndex = (colorIndex + 1) % colors.length;
        body.style.background = `linear-gradient(135deg, ${colors[colorIndex]} 0%, ${colors[(colorIndex + 1) % colors.length]} 100%)`;
    }, 5000);
    
    // Add a simple animation to any containers
    const containers = document.querySelectorAll('.container');
    containers.forEach(function(container) {
        container.style.transition = 'transform 0.3s ease';
        
        container.addEventListener('mouseenter', function() {
            this.style.transform = 'scale(1.02)';
        });
        
        container.addEventListener('mouseleave', function() {
            this.style.transform = 'scale(1)';
        });
    });
    
    // Log some information about the static file serving
    console.log('This JavaScript file is being served by the rtr static file handler!');
    console.log('Features:');
    console.log('- Automatic Content-Type detection');
    console.log('- Security protections');
    console.log('- 404 handling for missing files');
    console.log('- Seamless integration with rtr routing');
});

// Utility function to make API calls to our static data
function loadStaticData() {
    fetch('/static/data.json')
        .then(response => {
            if (!response.ok) {
                throw new Error('Data file not found');
            }
            return response.json();
        })
        .then(data => {
            console.log('Static data loaded:', data);
        })
        .catch(error => {
            console.log('Error loading static data:', error.message);
        });
}
