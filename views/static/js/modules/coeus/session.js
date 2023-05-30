// Handle session start
export function handleStartSession(input) {
    const courseRow = document.querySelector(`.course-row[data-sectionid="${input.sectionID}"]`);
    courseRow.classList.add("my-course-card-active-true");
    courseRow.classList.remove("my-course-card-active-false");

    const primaryFonts = courseRow.querySelectorAll('.my-course-card-primary-font-active-false');
    primaryFonts.forEach(font => {
        font.classList.add("my-course-card-primary-font-active-true");
        font.classList.remove("my-course-card-primary-font-active-false");
    });

    const secondaryFonts = courseRow.querySelectorAll('.my-course-card-secondary-font-active-false');
    secondaryFonts.forEach(font => {
        font.classList.add("my-course-card-secondary-font-active-true");
        font.classList.remove("my-course-card-secondary-font-active-false");
    });

    const tertiaryFonts = courseRow.querySelectorAll('.my-course-card-tertiary-font-active-false');
    tertiaryFonts.forEach(font => {
        font.classList.add("my-course-card-tertiary-font-active-true");
        font.classList.remove("my-course-card-tertiary-font-active-false");
    });

    const buttons = courseRow.querySelectorAll('.my-course-card-button-active-false');
    buttons.forEach(button => {
        button.classList.add("my-course-card-button-active-true");
        button.classList.remove("my-course-card-button-active-false");
    });
}

// Handle session end
export function handleEndSession(input) {
    const courseRow = document.querySelector(`.course-row[data-sectionid="${input.sectionID}"]`);
    if (courseRow) {
        courseRow.classList.add("my-course-card-active-false");
        courseRow.classList.remove("my-course-card-active-true");

        const primaryFonts = courseRow.querySelectorAll('.my-course-card-primary-font-active-true');
        primaryFonts.forEach(font => {
            font.classList.add("my-course-card-primary-font-active-false");
            font.classList.remove("my-course-card-primary-font-active-true");
        });

        const secondaryFonts = courseRow.querySelectorAll('.my-course-card-secondary-font-active-true');
        secondaryFonts.forEach(font => {
            font.classList.add("my-course-card-secondary-font-active-false");
            font.classList.remove("my-course-card-secondary-font-active-true");
        });

        const tertiaryFonts = courseRow.querySelectorAll('.my-course-card-tertiary-font-active-true');
        tertiaryFonts.forEach(font => {
            font.classList.add("my-course-card-tertiary-font-active-false");
            font.classList.remove("my-course-card-tertiary-font-active-true");
        });

        const buttons = courseRow.querySelectorAll('.my-course-card-button-active-true');
        buttons.forEach(button => {
            button.classList.add("my-course-card-button-active-false");
            button.classList.remove("my-course-card-button-active-true");
        });
    }
}




// Start class session
export function startClassSession(input) {
    // Placeholder for now
}

// End class session
export function endClassSession(input) {
    const sessionEndedModalElement = document.getElementById('end-session-modal');
    const sessionEndedModal = new mdb.Modal(sessionEndedModalElement);
    sessionEndedModal.show();

    countdown(5, function () {
        window.location.href = '/';
    });
}

// Handle participant joined
export function participantJoined(input) {
    const participantCount = document.getElementById("participant-count");
    const participants = input.count;
    participantCount.textContent = participants;
}

// Handle participant left
export function participantLeft() {
    const participantCount = document.getElementById("participant-count");
    let participants = parseInt(participantCount.textContent, 10);

    if (participants > 0) {
        participants--;
    }

    participantCount.textContent = participants;
}

// Update the votes of a question
export function updateVoteCount(input) {
    const questionID = input.questionID
    const votes = input.votes

    // Find the vote buttons with the specified questionID
    const voteButtons = document.querySelectorAll(`button.vote-up-btn[value="${questionID}"] p`);
  
    // Update the vote count and the arrow image for each button
    voteButtons.forEach(button => {
        button.innerHTML = `
            ${votes}
        `;
    });

    if (document.querySelector(`.unanswered-vote-count[data-vote-value="${questionID}"]`)) {
        const unansweredVotes = document.querySelectorAll(`.unanswered-vote-count[data-vote-value="${questionID}"]`);
        unansweredVotes.forEach(button => {
            button.textContent = `${votes} votes`;  
        });
    }
    

    // Rearrange the cards after updating the vote count
    rearrangeCards();
}

export function rearrangeCards() {
    const votesCardsContainer = document.querySelector('#votes-cards');
    const unansweredVotesCardsContainer = document.querySelector('#unanswered-votes-cards');

    if (unansweredVotesCardsContainer) {
        sortAndRearrangeCards(unansweredVotesCardsContainer, '.unanswered-vote-count');
    }

    sortAndRearrangeCards(votesCardsContainer, '.vote-up-btn p');
}

function sortAndRearrangeCards(container, voteCountSelector) {
    // Get all the cards within the container
    const cards = Array.from(container.querySelectorAll('.session-card-wrapper'));

    // Sort the cards based on the vote count
    const sortedCards = cards.sort((a, b) => {
        const voteCountA = parseInt(a.querySelector(voteCountSelector).textContent);
        const voteCountB = parseInt(b.querySelector(voteCountSelector).textContent);
        return voteCountB - voteCountA;
    });

    // Remove the current cards from the container
    container.innerHTML = '';

    // Append the sorted cards back to the container
    sortedCards.forEach(card => {
        container.appendChild(card);
    });
}
