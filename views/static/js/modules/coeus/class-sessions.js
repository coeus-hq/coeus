// Get tab elements and div elements
const createdCardsTab = document.getElementById('created-cards');
const votesCardsTab = document.getElementById('votes-cards');
const checkbox = document.getElementById('unanswered-toggle') ? document.getElementById('unanswered-toggle') : false;
const unansweredTopTab = document.getElementById('unanswered-votes-cards');
const unansweredNewestTab = document.getElementById('unanswered-created-cards');

// Toggle top tab
export function topTabToggle() {
    if (checkbox.checked) {
        unansweredTopTab.style.display = 'block';
        unansweredTopTab.classList.add('active', 'show');
        unansweredNewestTab.style.display = 'none';
        unansweredNewestTab.classList.remove('active', 'show');
    } else {
        votesCardsTab.style.display = 'block';
        createdCardsTab.style.display = 'none';
    }
}

// Toggle newest tab
export function newestTabToggle() {
    if (checkbox.checked) {
        unansweredNewestTab.style.display = 'block';
        unansweredNewestTab.classList.add('active', 'show');
        unansweredTopTab.style.display = 'none';
        unansweredTopTab.classList.remove('active', 'show');
    } else {
        createdCardsTab.style.display = 'block';
        votesCardsTab.style.display = 'none';
    }
}

// Show unanswered questions
export function showUnanswered() {
    if (votesCardsTab.classList.contains('active')) {
        unansweredTopTab.style.display = 'block';
        unansweredTopTab.classList.add('active', 'show');
        votesCardsTab.style.display = 'none';

    } else if (createdCardsTab.classList.contains('active')) {
        unansweredNewestTab.style.display = 'block';
        unansweredNewestTab.classList.add('active', 'show');
        createdCardsTab.style.display = 'none';
    }
}

// Hide unanswered questions
function hideUnanswered() {
    if (votesCardsTab.classList.contains('active')) {
        unansweredTopTab.style.display = 'none';
        unansweredTopTab.classList.remove('active', 'show');
        votesCardsTab.style.display = 'block';

    } else if (createdCardsTab.classList.contains('active')) {
        unansweredNewestTab.style.display = 'none';
        unansweredNewestTab.classList.remove('active', 'show');
        createdCardsTab.style.display = 'block';
    }
}

// Toggle unanswered questions checkbox
export function unansweredToggle() {
    const checkbox = document.getElementById('unanswered-toggle');

    const tabContent = document.getElementById('custom-tab-content');
    const tabPanes = tabContent.getElementsByClassName('tab-pane');

    for (let i = 0; i < tabPanes.length; i++) {
        if (tabPanes[i].classList.contains('active')) {
            if (checkbox.checked) {
                showUnanswered();
            } else {
                hideUnanswered();
            }
            break;
        }
    }
}

// Vote up function for questions
export function voteUp(event) {
    const button = event.currentTarget;
    const questionID = button.value;
    const userID = document.getElementById("userID").value;

    // Update image source attribute for both tabs
    const imgElements = document.querySelectorAll(`button[value="${questionID}"] img`);
    imgElements.forEach(imgElement => {
        imgElement.src = "/static/images/icon-arrow-up-true.svg";
    });

    // Send a POST request to server to upvote question
    fetch(`/api/vote-up/${questionID}`, {
        method: 'POST'
    })
        .then(response => {
            if (response.status == 200) {
            }
        })
        .catch(error => {
            console.log(error);
        });
};

// Mark question as answered
export function markQuestionAnswered(event) {
    const button = event.currentTarget;
    const questionID = button.value;

    // Send a POST request to server to mark question as answered
    fetch(`/api/mark-question/${questionID}`, {
        method: 'POST'
    })
        .then(response => {
            if (response.status == 200) {
                // Place holder for now
            }
        })
        .catch(error => {
            console.log(error);
        });
}

// Create question card
function createQuestionCard(question, moderatorType) {

    const formattedTime = formatTime24To12(question.CreatedAt);
    const card = document.createElement('div');
    card.className = 'card session-card-wrapper';
    card.setAttribute('data-question-id', question.ID);
    card.innerHTML = `
    <div class="question-card-time">
      ${formattedTime}
    </div>
    <div class="card-body dark-card-body">
    <strong class="session-q-text">Q.</strong>
    <p class="card-text session-card-text dark-card-text">
        ${question.Text}
      </p>
    </div>
    <div class="card-footer session-card-footer dark-card-footer">
      <button onclick="voteUp(event)" value="${question.ID}" class="vote-up-btn">
        <img src="/static/images/icon-arrow-up-${question.UserHasVoted}.svg" alt="">
        <p class="m-0 session-votes-font">${question.Votes}</p>
      </button>
      ${moderatorType == "student" ? '' : question.Answered == true ? '' : `<button onclick="markQuestionAnswered(event)" value="${question.ID}" class="mark-answered-btn">Answered?</button>`}
      <div class="session-answered-text">Answered:
        ${question.Answered ? '<span class="answered-true"> Yes </span>' : `<span class="answered-false" data-question-id="${question.ID}" > No </span>`}
      </div>
    </div>`;
    return card;
}

// Create unanswered question card 
function createUnansweredQuestionCard(question) {
    const formattedTime = formatTime24To12(question.CreatedAt);
    const card = document.createElement('div');
    card.className = 'card session-card-wrapper';
    card.setAttribute('data-question-id', question.ID);
    card.innerHTML = `
    <div class="question-card-time">
    ${formattedTime}
  </div>
  <div class="card-body dark-card-body">
  <strong class="session-q-text">Q.</strong>
  <p class="card-text session-card-text dark-card-text">
      ${question.Text}
    </p>
  </div>
  <div class="card-footer session-card-footer dark-card-footer">
      <p class="m-0 unanswered-vote-count session-votes-font" data-vote-value="${question.ID}">${question.Votes} votes</p>
    ${`<button onclick="markQuestionAnswered(event)" value="${question.ID}" class="mark-answered-btn">Mark Answered</button>`}
  </div>`;
    return card;
}

// Update questions
export function updateQuestions() {
    const url = window.location.href;
    const urlParts = url.split('/');
    const classSessionID = urlParts[urlParts.length - 2];

    fetch(`/api/questions/${classSessionID}`, {
        method: 'GET'
    })
        .then(response => response.json())
        .then(data => {

            const containerTime = document.getElementById('created-cards');
            const containerVote = document.getElementById('votes-cards');
            const containerTimeUnanswered = document.getElementById('unanswered-created-cards');
            const containerVoteUnanswered = document.getElementById('unanswered-votes-cards');

            containerTime.innerHTML = '';
            containerVote.innerHTML = '';
            containerTimeUnanswered.innerHTML = '';
            containerVoteUnanswered.innerHTML = '';

            data.questionsByTime.forEach(question => {
                const card = createQuestionCard(question, data.moderatorType);
                containerTime.appendChild(card);
            });

            data.questionsByVote.forEach(question => {
                const card = createQuestionCard(question, data.moderatorType);
                containerVote.appendChild(card);
            });

            data.questionsByVoteUnanswered.forEach(question => {
                const card = createUnansweredQuestionCard(question);
                containerVoteUnanswered.appendChild(card);
            });

            data.questionsByTimeUnanswered.forEach(question => {
                const card = createUnansweredQuestionCard(question);
                containerTimeUnanswered.appendChild(card);
            });
        })
        .catch(error => {
            // place holder for error handling right now
        });
}

// Check if page has questions container and call updateQuestions function
export function tryUpdateQuestions() {
    const containerTime = document.getElementById('created-cards');
    const containerVote = document.getElementById('votes-cards');

    if (containerTime && containerVote) {
        updateQuestions();
    }
}

tryUpdateQuestions();