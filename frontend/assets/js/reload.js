document.addEventListener("DOMContentLoaded", function() {

    connectToReload();
});

const reloadInterval = 500;
let version = null;
let reconnectTimeoutId = null;
// const maxRetryDelay = 10000;

function scheduleReconnect() {
	if (reconnectTimeoutId !== null) {
		return; // A reconnect is already scheduled
	}

	reconnectTimeoutId = setTimeout(function() {
		reconnectTimeoutId = null;
		connectToReload();
	}, reloadInterval);

	// Exponential backoff with a cap
	// currentRetryDelay = Math.min(maxRetryDelay, Math.floor(currentRetryDelay * 2));
}

function connectToReload() {
    const wsConnection = new WebSocket( "ws://" + window.location.host + "/reload");

    wsConnection.onmessage = function(event) {
        if (version === null) {
            console.log("[hot reload] server version: ", event.data);
            version = event.data;
        } else if (version !== event.data) {
            console.log("[hot reload] server version changed, reloading...");
            location.reload();
        }
    };

    wsConnection.onopen = function() {
        console.log("[hot reload] server connected");
		// Reset backoff when we successfully connect
		currentRetryDelay = reloadInterval;
		if (reconnectTimeoutId !== null) {
			clearTimeout(reconnectTimeoutId);
			reconnectTimeoutId = null;
		}
    };

	wsConnection.onerror = function(event) {
		// Do not schedule here exclusively; close may also fire. Use guarded scheduler.
		scheduleReconnect();
	};

	wsConnection.onclose = function(event) {
		scheduleReconnect();
	};
}