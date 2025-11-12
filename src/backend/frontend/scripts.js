document.getElementById('helloBtn').addEventListener('click', async () => {
  const reply = document.getElementById('reply');
  reply.textContent = 'Thinking…';
  try {
    const res = await fetch('/api/hello');
    const data = await res.json();
    reply.textContent = data.message;
  } catch (err) {
    reply.textContent = 'Error: ' + err.message;
  }
});