// Get elements
const hideChatBtn = document.getElementById("hide-chat-btn");
const questionTextarea = document.getElementById("questionTextarea");
const chatbox = document.getElementById("chatbox");
const bottomNavbar = document.getElementById("bottom-navbar");
const askBtn = document.querySelector('#ask-a-question-btn');

// Show chatbox
export function showChatbox() {
  askBtn.style.display = 'none';
  chatbox.style.display = 'block';
}

// Hide chatbox
export function hideChat() {
  chatbox.style.display = "none";
  askBtn.style.display = "flex";
}

// Textarea character counter
export function limitTextarea(textarea, maxChars) {
  if (textarea.value.length > maxChars) {
    textarea.value = textarea.value.substring(0, maxChars);
  }
  const charCountSpan = document.getElementById("charCount");
  charCountSpan.innerHTML = `${textarea.value.length}/${maxChars}`;
  if (textarea.value.length === maxChars) {
    charCountSpan.style.color = "red";
  } else {
    charCountSpan.style.color = "black";
  }
}

// Submit chatbox form
export function submitChatbox(event, classSessionId) {
  event.preventDefault();
  const questionText = document.getElementById("questionTextarea").value.trim();
  const userID = document.getElementById("userID").value;
  const formData = new FormData();
  formData.append("questionText", questionText);

  fetch(`/api/questions/${classSessionId}`, {
    method: "POST",
    headers: {},
    body: formData
  }).then(function (response) {
    if (response.status == 200) {
      const questionPostedMessage = document.getElementById("question-posted-message");

      document.getElementById("questionTextarea").value = '';
      const charCountSpan = document.getElementById("charCount");
      charCountSpan.innerHTML = "0/140";
      charCountSpan.style.color = "black";

      hideChat()
    } else {
      const questionNotPostedMessage = document.getElementById("question-not-posted-message");
      questionNotPostedMessage.style.display = "block";
      setTimeout(function () {
        questionNotPostedMessage.style.display = "none";
      }, 3000);
    }
  });
}

// Prevent empty questions from being submitted
export function preventEmptyQuestion(event, form) {
  const questionTextarea = document.getElementById("questionTextarea");
  const hiddenInputClassSessionID = document.getElementById("classSessionID").value;
  if (questionTextarea.value == "") {
    const enterAQuestionMessage = document.getElementById("enter-a-question");
    enterAQuestionMessage.style.display = "block";
    setTimeout(function () {
      enterAQuestionMessage.style.display = "none";
    }, 3000);
    event.preventDefault();
  } else {
    submitChatbox(event, hiddenInputClassSessionID);
  }
}