// renderQuestionCard renders a question card with the given question data
function renderQuestionCard(question, moderatorStatus, timezone) {

  const formattedTime = formatTime24To12(question.createdAt, timezone);
  const card = document.createElement('div');
  card.className = 'card session-card-wrapper';
  card.setAttribute('data-question-id', question.questionID);
  card.innerHTML = `
  <div class="question-card-time">
  ${formattedTime}
</div>
<div class="card-body dark-card-body">
<strong class="session-q-text">Q.</strong>
<p class="card-text session-card-text dark-card-text">
    ${question.text}
  </p>
</div>
<div class="card-footer session-card-footer dark-card-footer ">
  <button onclick="voteUp(event)" value="${question.questionID}" class="vote-up-btn">
    <img src="/static/images/icon-arrow-up-${question.answered}.svg" alt="">
    <p class="m-0 session-votes-font">${question.votes}</p>
  </button>
  ${moderatorStatus == "student" ? '' : question.answered == true ? '' : `<button onclick="markQuestionAnswered(event)" value="${question.questionID}" class="mark-answered-btn">Answered?</button>`}
  <div class="session-answered-text">Answered:
    ${question.answered ? '<span class="answered-true"> Yes </span>' : `<span class="answered-false" data-question-id="${question.questionID}" > No </span>`}
  </div>
</div>`;
  return card;
}

function renderUnansweredQuestionCard(question, timezone) {
  const formattedTime = formatTime24To12(question.createdAt, timezone);
  const card = document.createElement('div');
  card.className = 'card session-card-wrapper';
  card.setAttribute('data-question-id', question.questionID);
  card.innerHTML = `
  <div class="question-card-time">
  ${formattedTime}
</div>
<div class="card-body dark-card-body">
<strong class="session-q-text">Q.</strong>
  <p class="card-text session-card-text dark-card-text">
    ${question.text}
  </p>
</div>
<div class="card-footer session-card-footer dark-card-footer">
    <p class="m-0 unanswered-vote-count session-votes-font" data-vote-value="${question.questionID}">${question.votes} votes</p>
  ${`<button onclick="markQuestionAnswered(event)" value="${question.questionID}" class="mark-answered-btn">Mark Answered</button>`}
</div>`;
  return card;
}

// Render a new question card in the UI
export function renderNewQuestion(question) {
  const welcomeAlert = document.getElementById('welcome-alert');
  if (welcomeAlert) {
    welcomeAlert.remove();
  }

  const moderatorStatus = document.getElementById('moderator-status').value;
  const containerTime = document.getElementById('created-cards');
  const containerVote = document.getElementById('votes-cards');
  const containerTimeUnanswered = document.getElementById('unanswered-created-cards');
  const containerVoteUnanswered = document.getElementById('unanswered-votes-cards');
  const timezone = document.getElementById('timezone').value;

  const cardTime = renderQuestionCard(question, moderatorStatus, timezone);
  const cardVote = renderQuestionCard(question, moderatorStatus, timezone);

  const cardTimeUnanswered = renderUnansweredQuestionCard(question, timezone);
  const cardVoteUnanswered = renderUnansweredQuestionCard(question, timezone);

  if (containerTime.firstChild) {
    containerTime.insertBefore(cardTime, containerTime.firstChild);
    containerTimeUnanswered.insertBefore(cardTimeUnanswered, containerTimeUnanswered.firstChild);
  } else {
    containerTime.appendChild(cardTime);
    containerTimeUnanswered.appendChild(cardTimeUnanswered);
  }
  containerVote.appendChild(cardVote);
  containerVoteUnanswered.appendChild(cardVoteUnanswered);
}


// Remove the 'answered' button and update the 'answered' status in the UI
export function removeAnsweredBtn(question) {
  const markAnsweredBtns = document.querySelectorAll(`.mark-answered-btn[value="${question.questionID}"]`);
  if (markAnsweredBtns) {
    markAnsweredBtns.forEach(btn => btn.style.display = "none");
  }

  const answeredStatuses = document.querySelectorAll(`.answered-false[data-question-id="${question.questionID}"]`);
  if (answeredStatuses) {
    answeredStatuses.forEach(status => {
      status.innerHTML = "Yes";
      status.classList.remove("answered-false");
      status.classList.add("answered-true");
    });
  }
}

// Remove the question card from the unanswered questions tab when it is marked as answered
export function removeCard(question) {
  const unansweredTopTab = document.getElementById('unanswered-votes-cards');
  const unansweredNewestTab = document.getElementById('unanswered-created-cards');

  const topTabCard = unansweredTopTab.querySelector(`.session-card-wrapper[data-question-id="${question.questionID}"]`);
  const newestTabCard = unansweredNewestTab.querySelector(`.session-card-wrapper[data-question-id="${question.questionID}"]`);

  // Add fade-out class to the cards
  topTabCard.classList.add('fade-out');
  newestTabCard.classList.add('fade-out');

  // Remove the cards after the fade-out animation is complete
  setTimeout(() => {
    topTabCard.remove();
    newestTabCard.remove();
  }, 300);
}


// Show a welcome message if there are no questions
export function checkEmptyCards() {
  const createdCards = document.getElementById("created-cards");

  if (!createdCards.innerHTML.trim()) {
    const message = `
        <div id="welcome-alert" class="alert alert-info" role="alert">
          Welcome! It's time to ask your first question!
        </div>
      `;

    createdCards.innerHTML = message;
  }
}

// Call the checkEmptyCards function when the document is ready and the URL contains 'class-session'
document.addEventListener("DOMContentLoaded", function () {
  if (window.location.href.indexOf("class-session") !== -1) {
    // Wait for a moment before calling the checkEmptyCards function
    setTimeout(checkEmptyCards, 250);
  }
});