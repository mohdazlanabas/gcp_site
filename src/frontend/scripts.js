// API Button Handler
document.getElementById('helloBtn').addEventListener('click', async () => {
  const reply = document.getElementById('reply');
  const button = document.getElementById('helloBtn');

  // Add loading state
  button.classList.add('loading');
  button.disabled = true;
  reply.textContent = 'Thinking…';

  try {
    const res = await fetch('/api/hello');
    const data = await res.json();
    reply.textContent = data.message;
  } catch (err) {
    reply.textContent = '❌ Error: ' + err.message;
  } finally {
    button.classList.remove('loading');
    button.disabled = false;
  }
});

// Display Current Time
function updateVisitTime() {
  const timeElement = document.getElementById('visitTime');
  const now = new Date();

  const options = {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: true
  };

  const formattedTime = now.toLocaleString('en-US', options);
  timeElement.textContent = formattedTime;
}

// Update time immediately and then every second
updateVisitTime();
setInterval(updateVisitTime, 1000);

// Detect Visitor Location using IP Geolocation
async function detectLocation() {
  const locationElement = document.getElementById('visitLocation');

  try {
    // Try using ipapi.co (free, no API key needed)
    const response = await fetch('https://ipapi.co/json/');

    if (!response.ok) {
      throw new Error('Location service unavailable');
    }

    const data = await response.json();

    // Format location string
    let locationString = '';

    if (data.city && data.country_name) {
      locationString = `${data.city}, ${data.country_name}`;
    } else if (data.country_name) {
      locationString = data.country_name;
    } else if (data.country_code) {
      locationString = data.country_code;
    } else {
      locationString = 'Location unknown';
    }

    // Add additional info if available
    if (data.region) {
      locationString = `${data.city}, ${data.region}, ${data.country_name}`;
    }

    // Add flag emoji if country code is available
    if (data.country_code) {
      const flagEmoji = getFlagEmoji(data.country_code);
      locationString = `${flagEmoji} ${locationString}`;
    }

    locationElement.textContent = locationString;

    // Store additional data as tooltip
    if (data.timezone) {
      locationElement.title = `Timezone: ${data.timezone}\nIP: ${data.ip || 'Unknown'}`;
    }

  } catch (error) {
    console.error('Location detection error:', error);

    // Fallback: Try to get timezone-based location
    try {
      const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
      locationElement.textContent = `Timezone: ${timezone}`;
    } catch (tzError) {
      locationElement.textContent = 'Location unavailable';
    }
  }
}

// Convert country code to flag emoji
function getFlagEmoji(countryCode) {
  const codePoints = countryCode
    .toUpperCase()
    .split('')
    .map(char => 127397 + char.charCodeAt());
  return String.fromCodePoint(...codePoints);
}

// Initialize location detection
detectLocation();

// Add smooth scroll behavior for any future internal links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
  anchor.addEventListener('click', function (e) {
    e.preventDefault();
    const target = document.querySelector(this.getAttribute('href'));
    if (target) {
      target.scrollIntoView({
        behavior: 'smooth',
        block: 'start'
      });
    }
  });
});

// Log app info to console
console.log('%c🚀 Google Cloud App Engine', 'font-size: 20px; font-weight: bold; color: #2563eb;');
console.log('%cPowered by Go 1.22', 'font-size: 14px; color: #64748b;');
console.log('%cServed from: Singapore (asia-southeast2)', 'font-size: 14px; color: #64748b;');
console.log('%cInstance: F1 Class', 'font-size: 14px; color: #64748b;');
