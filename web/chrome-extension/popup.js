document.getElementById("checkBtn").onclick = async () => {
  const url = document.getElementById("url").value.trim();
  const resultBox = document.getElementById("result");

  if (!url) {
    resultBox.textContent = "Please enter a URL";
    return;
  }

  resultBox.textContent = "Checking...";

  try {
    const res = await fetch(`http://localhost:8080/api/v1/analyze?url=${encodeURIComponent(url)}`);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = await res.json();

    const verdict = data.result?.verdict ?? "Unknown";
    const finalScore = data.result?.final_score ?? "N/A";

    // resultBox.textContent = `Verdict: ${verdict}\n , Trust Score: ${finalScore}`;
    // resultBox.textContent = `Verdict: ${verdict}\n , Trust Score: ${finalScore}`;

    resultBox.innerHTML = `
      Verdict: <strong>${verdict}</strong><br>
      Trust Score: <strong>${finalScore}</strong><br>
      <a href="#" id="moreDetails">Click here to see more details</a>
    `;

    document.getElementById("moreDetails").onclick = (e) => {
      e.preventDefault();
      // Open the full extension popup or a new tab with details
      chrome.tabs.create({
        url: `http://localhost:5173/?q=${encodeURIComponent(url)}`
      });
    };

  } catch (err) {
    resultBox.textContent = "Error: " + err.message;
  }
};
