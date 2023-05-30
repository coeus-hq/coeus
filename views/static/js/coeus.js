import * as signinCreateAccount from './modules/coeus/signin-create-account.js';
import * as myCourses from './modules/coeus/my-courses.js';
import * as classSessions from './modules/coeus/class-sessions.js';
import * as session from './modules/coeus/session.js';
import * as questions from './modules/coeus/questions.js';
import * as chatbox from './modules/coeus/chatbox.js';
import * as settings from './modules/coeus/settings.js';
import * as classSections from './modules/coeus/class-sections.js';
import * as utils from './modules/coeus/utils.js';
import * as passwordReset from './modules/coeus/password-reset.js';
import './modules/coeus/websockets.js';

// Expose all imported functions to the global scope
Object.assign(window, {
  ...signinCreateAccount,
  ...myCourses,
  ...classSessions,
  ...session,
  ...questions,
  ...chatbox,
  ...settings,
  ...classSections,
  ...utils,
  ...passwordReset
});