const CACHE_NAME = "clock-app-cache-v1";
const URLS_TO_CACHE = [
	"/",
	"/about",
	"/style.css",
	"/script.js",
	"/manifest.json",
	"/icons/icon-192.png",
	"/icons/icon-512.png",
];

self.addEventListener("install", (event) => {
	event.waitUntil(
		caches.open(CACHE_NAME).then((cache) => cache.addAll(URLS_TO_CACHE))
	);
});

self.addEventListener("fetch", (event) => {
	if (event.request.url.includes("/api/time")) {
		event.respondWith(
			fetch(event.request)
				.then((res) => {
					const clone = res.clone();
					caches.open(CACHE_NAME).then((cache) => cache.put(event.request, clone));
					return res;
				})
				.catch(() => caches.match(event.request))
		);
		return;
	}

	event.respondWith(
		caches.match(event.request).then((cached) => cached || fetch(event.request))
	);
});

self.addEventListener("activate", (event) => {
	event.waitUntil(
		caches.keys().then((keys) =>
			Promise.all(keys.filter((k) => k !== CACHE_NAME).map((k) => caches.delete(k)))
		)
	);});
