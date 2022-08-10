var pwd = document.getElementById("password");
var confirmPwd = document.getElementById("confirm_password");
var letter = document.getElementById("letterCheckFeedback");
var capital = document.getElementById("capitalCheckFeedback");
var number = document.getElementById("numberCheckFeedback");
var length = document.getElementById("lengthCheckFeedback");
var match = document.getElementById("matchCheckFeedback");
var submit = document.getElementById("submit");
var form = document.getElementById("profileForm");
var messageBox = document.getElementById("message");
var specialChars = "!@#$%^&*()-_=+[{]}\\|;:'\",<.>/?`~";  // unused at moment

// When the user clicks on the password field, show the message box
pwd.onfocus = function() {
  document.getElementById("message").style.display = "block";
}

// When the user clicks outside of the password field, hide the message box
pwd.onblur = function() {
  //document.getElementById("message").style.display = "none";
}
confirmPwd.onfocus = function() {
  document.getElementById("message").style.display = "block";
}

// When the user clicks outside of the password field, hide the message box
confirmPwd.onblur = function() {
  //document.getElementById("message").style.display = "none";
}

// When the user starts to type something inside the password field
var lowerCaseLettersCheck = false;
var upperCaseLettersCheck = false;
var numbersCheck = false;
var lengthCheck = false;
var matchCheck = false;

pwd.onkeyup = function() {
  // Validate lowercase letters
  var lowerCaseLetters = /[a-z]/g;
  if(pwd.value.match(lowerCaseLetters)) {  
    letter.classList.remove("invalid");
    letter.classList.add("valid");
    lowerCaseLettersCheck = true;
  } else {
    letter.classList.remove("valid");
    letter.classList.add("invalid");
    lowerCaseLettersCheck = false;
  }
  
  // Validate capital letters
  var upperCaseLetters = /[A-Z]/g;
  if(pwd.value.match(upperCaseLetters)) {  
    capital.classList.remove("invalid");
    capital.classList.add("valid");
    upperCaseLettersCheck = true;
  } else {
    capital.classList.remove("valid");
    capital.classList.add("invalid");
    upperCaseLettersCheck = false;
  }

  // Validate numbers
  var numbers = /[0-9]/g;
  if(pwd.value.match(numbers)) {  
    number.classList.remove("invalid");
    number.classList.add("valid");
    numbersCheck = true;
  } else {
    number.classList.remove("valid");
    number.classList.add("invalid");
    numbersCheck = false;
  }
  
  // Validate length
  if(pwd.value.length >= 8) {
    length.classList.remove("invalid");
    length.classList.add("valid");
    lengthCheck = true;
  } else {
    length.classList.remove("valid");
    length.classList.add("invalid");
    lengthCheck = false;
  }
}

confirmPwd.onkeyup = function() {
    if (pwd.value == confirmPwd.value) {
        match.classList.remove("invalid");
        match.classList.add("valid");
        matchCheck = true;
    } else {
        match.classList.remove("valid");
        match.classList.add("invalid");
        matchCheck = false;
    }

    // color message box
    if (lowerCaseLettersCheck && upperCaseLettersCheck && numbersCheck && lengthCheck && matchCheck) {
        messageBox.classList.remove("invalid-box");
        messageBox.classList.add("valid-box");
    } else {
        messageBox.classList.remove("valid-box");
        messageBox.classList.add("invalid-box");
    }
}
