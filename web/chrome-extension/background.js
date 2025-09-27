chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === "analyze_url") {
    fetch(`http://localhost:8080/api/v1/analyze?url=${encodeURIComponent(message.url)}`)
      .then(res => res.json())
      .then(data => sendResponse({ success: true, data }))
      .catch(err => sendResponse({ success: false, error: err.message }));

    return true; // required for async sendResponse
  }
});
