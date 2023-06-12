[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Official Website](https://img.shields.io/badge/Website-coeus.education-blue)](https://coeus.education/)
[![Deployed App](https://img.shields.io/badge/Deployed_App-Visit-blue.svg)](http://dev.coeus.education/)
![Open Source](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)



![Coeus banner](https://coeus.education/images/coeus-banner.png)

Coeus Education is an open-source project that aspires to enhance communication between students and instructors. Designed for inclusivity, our platform allows students to ask questions anonymously, thereby fostering active participation, curiosity, and a judgement-free learning environment.

In recognizing that many students share common queries, Coeus aims to ensure every question is addressed. Our user-friendly interface simplifies the process for instructors to manage courses and sessions, while students can engage easily, even taking on roles such as teacher assisants.

As a community-driven initiative, we invite contributions from developers, educators, and students.

Join us to transform the learning experience with Coeus Education.

## Technologies

Coeus Education leverages the following stack:

- **Backend**: Golang and the Gin web framework
- **Database**: SQLite3 for data management
- **Frontend**: Go's template engine with HTML, JavaScript, and Bootstrap 5
- **Email Service**: SendGrid for user communication

These technologies combine to create Coeus.


## Getting Started

### Follow these steps to get started with Coeus:

**Requirements**: Ensure that you have the following requirements installed on your system:
- Go programming language: https://golang.org/dl/

#### Mac

**1. Clone the Repository**: Open your terminal or command prompt and clone the Coeus repository by running the following command:
```bash
git clone https://github.com/coeus-hq/coeus.git
```

**2. Build and Run the Application**: Navigate to the project's root directory and run the following command to build and run the Coeus server:

```bash
go run .
```
The server should start, and you'll be able to access Coeus at http://localhost:8080.

#### Windows

**1. Install GCC**:
   The Go package that we're using (go-sqlite3) requires CGO, which in turn needs GCC to work. The easiest way to get GCC on your Windows machine is through the MinGW-w64 project.

   a. Go to the [MinGW-w64 project files page on SourceForge](https://sourceforge.net/projects/mingw-w64/files/mingw-w64/).

   b. Click on the "MinGW-W64 GCC-8.1.0" link. You'll see several versions of MinGW-w64 for different systems. Choose the one that corresponds to your system (typically, you'll want `x86_64-posix-seh`).

![README COEUS WINDOWS HELP SS](https://github.com/asimbaig95/coeus/assets/88279562/44823119-6290-496f-b89a-ba6fd8e21685)


   c. This will download a .7z archive file. To extract this, you will need a tool capable of extracting .7z files, such as [7-Zip](https://www.7-zip.org/download.html). Once 7-Zip is installed, right-click the .7z file you downloaded, navigate to 7-Zip in the context menu and click on "Extract here".

   d. You should now have a MinGW-w64 directory. Move this directory to a safe location (e.g., `C:\\`). Take note of the path to the `bin` directory inside this MinGW-w64 directory, you'll need it for the next step.

   e. Add the path to the `bin` directory to your PATH environment variable. You can do this by searching for "Environment Variables" in the Start Menu, clicking on "Edit the system environment variables", clicking on the "Environment Variables" button, scrolling to "Path" under "System variables", clicking "Edit", and then "New". Paste the path to the `bin` directory and click OK on all windows to close them.
   
   You path will likely look like: 
```bash
   C:\mingw64\bin
```
**1.5. Confirm Installation**:
   
   After you've installed MinGW-w64 and added it to your PATH environment variable, you should confirm that the installation was successful. Open a new Command Prompt window (you can do this by typing "cmd" into the Start Menu and pressing Enter) and type the following command:

   ```bash
   gcc --version
```

**2. Clone the Repository**: Open your terminal or command prompt and clone the Coeus repository by running the following command:
```bash
git clone https://github.com/coeus-hq/coeus.git
```

**2. Build and Run the Application**: Navigate to the project's root directory and run the following command to build and run the Coeus server:

```bash
go run .
```
The server should start, and you'll be able to access Coeus at http://localhost:8080.


**3. Set Up Admin Account**: Upon first run, Coeus will be seeded with testing data, but you have the option to reset and reseed the database in the organization settings while logged in as an admin, reseting the database will trigger an onboarding wizard. Follow the on-screen instructions to create an admin account with the necessary credentials.

**4. Explore Coeus**: Once the server is up and running, you can start using Coeus. 
you may,

Login as:
student@coeus.education
instructor@coeus.education
admin@coeus.education

Password:
coeus

That's it! You have successfully set up Coeus and are ready to try it out and test it.

## Trouble Shooting

Issues may occur with cookies when switching between user accounts in Coeus.

### How to fix:
1. **Clear browser's cookies**
2. **Restart browser**
3. **Try different browser or device**

If issues persist, please contact a repo manager.

## Running Model Tests

To run tests for the models, navigate to the `models` directory and execute the following command:

```bash
cd models
go test -v
```
This will run all tests and display verbose output.

#### ⚠️ Note that the model tests depend on the default seed data ⚠️



## Meet the Team

## [Asim](https://github.com/asimbaig95)
The visionary founder of Coeus. With his expertise in software architecture and deployment, he has laid the foundation for Coeus to deliver a robust and scalable platform.

**Role:** 
- Founder 💡
- Software Architect 🏗️
- Deployment Specialist ⚙️

## [Collin](https://github.com/Collin-W)
The full stack developer behind Coeus. Built the front and back-end of Coeus.

**Role:**
- Full Stack Developer 🖥️
- Repo Manager (Code/Pull requests related) 📁

## [Molly](https://github.com/mollyshwiff)
The creative force, project manager, and tester behind Coeus. Crafting engaging experiences that prioritize usability.

**Role:**
- UI/UX Designer 🎨
- Project Manager 📋
- Repo Manager (UX/UI related) 📁
- Tester 🔍

## Your Name Here
 If you have a desire to contribute to open source software, and are eager to learn, we would love to hear from you!

**Role:**
- Making Coeus a great open source educational software 🚀💪


Please feel free to reach out to our team with any inquiries or feedback.

## Contributing

We welcome contributions from the community! Please read our [CONTRIBUTING.md](CONTRIBUTING.md) file for more information on how to get involved.


