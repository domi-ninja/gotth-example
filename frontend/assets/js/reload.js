document.addEventListener("DOMContentLoaded", function() {

    connectToReload();
});

const reloadInterval = 500;
let version = null;

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
    };

    wsConnection.onerror = function(event) {
        setTimeout(connectToReload, reloadInterval);
    };

    wsConnection.onclose = function(event) {
        setTimeout(connectToReload, reloadInterval);
    };
}