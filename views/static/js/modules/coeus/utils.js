// Displays countdown in the modal before redirecting to the home page
export function countdown(seconds, callback) {
  let remainingSeconds = seconds;

  const intervalId = setInterval(() => {
    remainingSeconds--;
    document.getElementById("end-count-down").textContent = remainingSeconds;

    if (remainingSeconds <= 0) {
      clearInterval(intervalId);
      callback();
    }
  }, 1000);
}

// This is a hack to get the back button to work in the course section page breadcrumbs
export function backButtonClick() {
  window.history.back();
}

// formatTime24To12 converts a 24-hour time string to a 12-hour time string with AM/PM
export function formatTime24To12(time24, timezone = null) {
  // Parse the hours and minutes from the input time string
  const [hours24, minutes] = time24.split(':');

  if (timezone !== null) {
    const totalMinutes = parseInt(hours24) * 60 + parseInt(minutes) + parseInt(timezone);

    // Convert the total minutes back to hours and minutes, adjusting for 24-hour time
    const adjustedHours24 = Math.floor(totalMinutes / 60) % 24;
    const adjustedMinutes = totalMinutes % 60;

    // Convert the adjusted 24-hour time to a 12-hour format
    const hours12 = ((adjustedHours24 % 12) || 12);
    const adjustedMinutesPadded = adjustedMinutes.toString().padStart(2, '0');
    const period = adjustedHours24 < 12 ? 'a.m.' : 'p.m.';

    return `${hours12}:${adjustedMinutesPadded} ${period}`;
  } else {
    const hours12 = ((hours24 % 12) || 12);
    const period = hours24 < 12 ? 'a.m.' : 'p.m.';

    return `${hours12}:${minutes} ${period}`;
  }
}
