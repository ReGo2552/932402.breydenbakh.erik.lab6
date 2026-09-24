const express = require("express");
const app = express();
const PORT = process.env.PORT || 3000;

app.get("/api/time", (req, res) => {
	try {
		const now = new Date();
		const use12h = req.query.format === "12";

		const time = now.toLocaleTimeString("en-GB", {
			hour12: use12h,
			hour: "2-digit",
			minute: "2-digit",
			second: "2-digit",
		});

		const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

		res.json({ time, timezone, format: use12h ? "12h" : "24h" });
	} catch (err) {
		res.status(500).json({ error: "Internal server error" });
	}
});

app.listen(PORT, () => {
	console.log(`Server started on port ${PORT}`);
});
