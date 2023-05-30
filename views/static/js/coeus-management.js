import * as courseModals from './modules/management/course-modals.js';
import * as courseTable from './modules/management/course-table.js';
import * as userModals from './modules/management/user-modals.js';
import * as userTable from './modules/management/user-table.js';
import * as organizationSettings from './modules/management/organization-settings.js';
import * as utils from './modules/management/utils.js';
import * as attendance from './modules/management/attendance.js';
import * as onboarding from './modules/management/onboarding.js';

// Expose all imported functions to the global scope
Object.assign(window, {
  ...courseModals,
  ...courseTable,
  ...userModals,
  ...userTable,
  ...organizationSettings,
  ...utils,
  ...attendance,
  ...onboarding
});