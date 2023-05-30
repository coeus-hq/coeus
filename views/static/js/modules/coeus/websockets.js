// Create a new WebSocket connection for general messages
const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
const host = location.host;
let generalWs = new WebSocket(`${protocol}//${host}/ws`);

generalWs.onopen = function () {
    // WebSocket connection opened
};

generalWs.onmessage = function (event) {
    // Parse the received JSON message
    const input = JSON.parse(event.data);

    // Check the action of the received input
    switch (input.action) {
        case "start-session":
            handleStartSession(input);
            break;
        case "end-session":
            handleEndSession(input);
            break;
        case "new-logo":
            // Update the logo image
            const logoImg = document.getElementById("organization-logo");
           
            logoImg.src = `${input.logoPath}`;
            break;
        default:
            console.log("Unknown action:", input.action);
    }
};

let classSessionID;
let classWs;

// If the class session page has a class-session-ID element, create a websocket connection
if (document.getElementById("class-session-ID")) {
    // Get the class session ID from the class-session-ID element
    classSessionID = document.getElementById("class-session-ID").value;
    classWs = new WebSocket(`${protocol}//${host}/ws/${classSessionID}`);

    classWs.onopen = function () {
        // WebSocket connection opened
    };

    classWs.onmessage = function (event) {
        // Parse the received JSON message
        const input = JSON.parse(event.data);

        // Check the action of the received input
        switch (input.action) {
            case "vote-up":
                updateVoteCount(input);
                break;
            case "new-question":
                renderNewQuestion(input);
                break;
            case "mark-question":
                removeAnsweredBtn(input);
                removeCard(input);
                break;
            case "start-session":
                startClassSession(input);
                break;
            case "end-session":
                endClassSession(input);
                break;
            case "participant-joined":
                participantJoined(input);
                break;
            case "participant-left":
                participantLeft();
                break;
            default:
                console.log("Unknown action:", input.action);
        }
    };

    classWs.onclose = function () {
        // WebSocket connection closed
    };

    classWs.onerror = function (error) {
        console.log("WebSocket error:", error);
    };
}