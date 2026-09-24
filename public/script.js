const timeDisplay = document.getElementById("time-display");
const timezoneDisplay = document.getElementById("timezone-display");
const formatToggle = document.getElementById("format-toggle");
const offlineNotice = document.getElementById("offline-notice");

async function updateTime() {
	const url = formatToggle.checked ? "/api/time?format=12" : "/api/time";

	try {
		const res = await fetch(url);
		if (!res.ok) throw new Error(`HTTP ${res.status}`);
		const data = await res.json();

		timeDisplay.textContent = data.time;
		timezoneDisplay.textContent = data.timezone;
		offlineNotice.hidden = true;
	} catch (err) {
		offlineNotice.hidden = false;
	}
}

if (timeDisplay && formatToggle) {
	updateTime();
	setInterval(updateTime, 1000);
	formatToggle.addEventListener("change", updateTime);
}

if ("serviceWorker" in navigator) {
	window.addEventListener("load", () => {
		navigator.serviceWorker.register("/sw.js");
	});
}
