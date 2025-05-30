// Get references to HTML elements
const loadingDisplay = document.getElementById('loading');
const errorDisplay = document.getElementById('error');
const matchDisplay = document.getElementById('match-display');
const teamNamesElement = document.getElementById('team-names');
const scoreElement = document.getElementById('score');
const statusElement = document.getElementById('status');
const lastEventElement = document.getElementById('last-event');
const simulateGoalBtn = document.getElementById('simulate-goal-btn');

// Configuration
const GRPC_PROXY_URL = 'http://localhost:8081'; // Your gRPC-Web proxy URL
const MATCH_ID = 'match-123'; // The ID of the match you want to track

// Initialize the gRPC client
// The client expects the gRPC-Web proxy URL
const grpcClient = new proto.match.MatchServiceClient(GRPC_PROXY_URL);

// Function to update the display with new match data
function updateMatchDisplay(match) {
    teamNamesElement.textContent = `${match.getHomeTeam() || "Home"} vs ${match.getAwayTeam() || "Away"}`;
    scoreElement.textContent = `${match.getHomeScore()} - ${match.getAwayScore()}`;
    statusElement.textContent = `Status: ${match.getStatus()}`;
    lastEventElement.textContent = `Last Event: ${match.getLastEvent() || "N/A"}`;

    loadingDisplay.style.display = 'none';
    errorDisplay.style.display = 'none';
    matchDisplay.style.display = 'block';
}

// Function to show error messages
function showErrorMessage(message) {
    loadingDisplay.style.display = 'none';
    matchDisplay.style.display = 'none';
    errorDisplay.style.display = 'block';
    errorDisplay.textContent = `Error: ${message}`;
    console.error(message);
}

// Function to fetch initial match data and then subscribe to updates
function setupMatchUpdatesStream() {
    console.log(`Attempting to get match updates for: ${MATCH_ID}`);

    const request = new proto.match.MatchRequest();
    request.setMatchId(MATCH_ID);

    // Call the streaming RPC method
    const stream = grpcClient.getMatchUpdates(request, {});

    stream.on('data', (response) => {
        // 'response' is an instance of proto.match.Match
        console.log('Received match update:', response.toObject());
        updateMatchDisplay(response);
    });

    stream.on('status', (status) => {
        // Handle stream status changes (e.g., connection status, error codes)
        console.log('Stream status:', status);
        if (status.code !== grpc.web.StatusCode.OK) {
            showErrorMessage(`Stream error: ${status.details} (Code: ${status.code})`);
        }
    });

    stream.on('end', () => {
        // Stream ended (e.g., server closed the stream)
        console.log('Match update stream ended.');
        // Potentially attempt to re-establish the stream or show a message
        showErrorMessage('Match update stream closed by server. Trying to reconnect...');
        // Simple reconnect logic (be careful with rapid reconnects in production)
        setTimeout(setupMatchUpdatesStream, 3000);
    });

    stream.on('error', (err) => {
        // Handle network errors or other stream-specific errors
        console.error('Stream error:', err);
        showErrorMessage(`Network or stream error: ${err.message}`);
        // Simple reconnect logic
        setTimeout(setupMatchUpdatesStream, 3000);
    });
}


// Event listener for the "Simulate Admin Goal" button

// Initialize on DOMContentLoaded
document.addEventListener('DOMContentLoaded', () => {
    // Start the streaming updates when the page loads
    setupMatchUpdatesStream();
});