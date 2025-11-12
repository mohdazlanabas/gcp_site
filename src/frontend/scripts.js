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

// Message Storage functionality
const messageInput = document.getElementById('messageInput');
const sendMessageBtn = document.getElementById('sendMessageBtn');
const refreshMessagesBtn = document.getElementById('refreshMessagesBtn');
const messagesContainer = document.getElementById('messagesContainer');

// Send message
async function sendMessage() {
  const message = messageInput.value.trim();

  if (!message) {
    alert('Please enter a message');
    return;
  }

  // Disable button during request
  sendMessageBtn.disabled = true;
  sendMessageBtn.classList.add('loading');

  try {
    const response = await fetch('/api/messages', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ message }),
    });

    const data = await response.json();

    if (data.success) {
      // Clear input
      messageInput.value = '';

      // Show success feedback
      messageInput.placeholder = '✓ Message sent! Type another...';
      setTimeout(() => {
        messageInput.placeholder = 'Type your message here...';
      }, 2000);

      // Refresh messages to show the new one
      await loadMessages();
    } else {
      alert('Error: ' + (data.error || 'Failed to send message'));
    }
  } catch (error) {
    console.error('Error sending message:', error);
    alert('Failed to send message. Please try again.');
  } finally {
    sendMessageBtn.disabled = false;
    sendMessageBtn.classList.remove('loading');
    messageInput.focus();
  }
}

// Load and display messages
async function loadMessages() {
  try {
    const response = await fetch('/api/messages');
    const data = await response.json();

    if (data.success) {
      displayMessages(data.messages);
    } else {
      console.error('Failed to load messages:', data.error);
    }
  } catch (error) {
    console.error('Error loading messages:', error);
  }
}

// Display messages in the UI
function displayMessages(messages) {
  if (!messages || messages.length === 0) {
    messagesContainer.innerHTML = '<div class="no-messages">No messages yet. Be the first to send one! 📝</div>';
    return;
  }

  // Build HTML for messages
  let html = '';
  messages.forEach((msg, index) => {
    const date = new Date(msg.timestamp);
    const timeAgo = getTimeAgo(date);
    const formattedTime = date.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });

    html += `
      <div class="message-item">
        <div class="message-header">
          <span class="message-label">Past Message ${index + 1}</span>
          <span class="message-time" title="${date.toLocaleString()}">${timeAgo}</span>
        </div>
        <div class="message-content">${escapeHtml(msg.content)}</div>
        <div class="message-id">ID: ${msg.id} • ${formattedTime}</div>
      </div>
    `;
  });

  messagesContainer.innerHTML = html;
}

// Calculate time ago
function getTimeAgo(date) {
  const seconds = Math.floor((new Date() - date) / 1000);

  if (seconds < 60) return 'Just now';
  if (seconds < 3600) return `${Math.floor(seconds / 60)} min ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)} hours ago`;
  return `${Math.floor(seconds / 86400)} days ago`;
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

// Event listeners for message storage
sendMessageBtn.addEventListener('click', sendMessage);

refreshMessagesBtn.addEventListener('click', async () => {
  refreshMessagesBtn.disabled = true;
  refreshMessagesBtn.classList.add('loading');
  await loadMessages();
  refreshMessagesBtn.disabled = false;
  refreshMessagesBtn.classList.remove('loading');
});

// Allow Enter key to send message
messageInput.addEventListener('keypress', (e) => {
  if (e.key === 'Enter') {
    sendMessage();
  }
});

// Load messages on page load
loadMessages();

// Log app info to console
console.log('%c🚀 Google Cloud App Engine', 'font-size: 20px; font-weight: bold; color: #2563eb;');
console.log('%cPowered by Go 1.22', 'font-size: 14px; color: #64748b;');
console.log('%cServed from: Singapore (asia-southeast2)', 'font-size: 14px; color: #64748b;');
console.log('%cInstance: F1 Class', 'font-size: 14px; color: #64748b;');
console.log('%c💾 Message Storage: In-Memory', 'font-size: 14px; color: #2563eb;');
