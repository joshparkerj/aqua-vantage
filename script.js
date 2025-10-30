const pageElements = {
  josephRiveraTitle: document.querySelector("main > h3"),
  friendlyOffer: document.querySelector("p.friendly-offer"),
  yourName: document.querySelector("label[for=your-name]"),
  addressLineOne: document.querySelector("label[for=address-line-1]"),
  addressLineTwo: document.querySelector("label[for=address-line-2]"),
  city: document.querySelector("label[for=address-city]"),
  state: document.querySelector("label[for=address-state]"),
  zipCode: document.querySelector("label[for=address-zip]"),
  phoneNumber: document.querySelector("label[for=phone-number]"),
  apptDateTime: document.querySelector("label[for=appt-date-time]"),
  languageToggle: document.querySelector("button.language-toggle"),

  formTitle: document.querySelector("form > h2"),
  cancellation: document.querySelector("form > p.cancellation"),
  emphaticRemark: document.querySelector("form > p.emphatic-remark"),
  briefPitch: document.querySelector("p.brief-pitch"),

  nextPage: document.querySelector("button.next-page"),
  previousPage: document.querySelector("button.previous-page"),
};

let language = "spanish";

const englishLanguageText = {
  josephRiveraTitle: "Water Analyst",
  friendlyOffer: "Job Opportunity",
  yourName: "Your Name:",
  addressLineOne: "Address Line 1:",
  addressLineTwo: "Address Line 2:",
  city: "City:",
  state: "State:",
  zipCode: "Zip Code:",
  phoneNumber: "Phone Number:",
  apptDateTime: "Appointment Date and Time:",
  languageToggle: "Español",
  formTitle: "Request your appointment for water analysis",
  cancellation:
    "If you must cancel your scheduled appointment, please call us as soon as possible.",
  emphaticRemark: "It's completely free and there's no obligation!",
  briefPitch:
    "You have the right to know the condition of your drinking water and how it may affect your HEALTH, your HOUSEHOLD BUDGET, and your PROPERTY.",
  nextPage: "Next Page",
  previousPage: "Previous Page",
};

const spanishLanguageText = {
  josephRiveraTitle: "Analista de Agua",
  friendlyOffer: "Oportunidad de Trabajo",
  yourName: "Nombre:",
  addressLineOne: "Dirección Línea 1:",
  addressLineTwo: "Dirección Línea 2:",
  city: "Ciudad:",
  state: "Estado:",
  zipCode: "Código Postal:",
  phoneNumber: "Número de Teléfono:",
  apptDateTime: "Fecha y Hora de la Cita:",
  languageToggle: "English",
  formTitle: "Confirmación de Cita para analizar su agua",
  cancellation:
    "Si Cree qus su cita no puede ser en este dia llamenos lo más pronto posible.",
  emphaticRemark: "ES GRATIS Y SIN NINGÚN COMPROMISO!",
  briefPitch:
    "Usted tiene el derecho de saber en qué condición le llega su agua potable; y como le puede afectar, su SALUD, su ECONOMIA DOMESTICA, y su PROPIEDAD.",
  nextPage: "Página Siguiente",
  previousPage: "Página Anterior",
};

pageElements.languageToggle.addEventListener("click", () => {
  let text;
  if (language === "spanish") {
    text = englishLanguageText;
    language = "english";
  } else if (language === "english") {
    text = spanishLanguageText;
    language = "spanish";
  } else {
    throw "you messed up the language 🤦";
  }

  Object.entries(pageElements).forEach(([elementName, element]) => {
    element.textContent = text[elementName];
  });
});

let currentPageNumber = 0;

const pages = [...document.querySelectorAll("[class*=page-number]")];

const pageTurnEventListener = (predicate, pageIncrement) => () => {
  if (predicate()) {
    currentPageNumber += pageIncrement;
    if (currentPageNumber === pages.length - 1) {
      pageElements.nextPage.classList.add("inactive");
    } else {
      pageElements.nextPage.classList.remove("inactive");
    }

    if (currentPageNumber === 0) {
      pageElements.previousPage.classList.add("inactive");
    } else {
      pageElements.previousPage.classList.remove("inactive");
    }
    pages.forEach((page, pageNumber) => {
      if (pageNumber < currentPageNumber) {
        page.classList.remove("on-screen");
        page.classList.remove("off-right");
        page.classList.add("off-left");
      } else if (pageNumber === currentPageNumber) {
        page.classList.remove("off-right");
        page.classList.remove("off-left");
        page.classList.add("on-screen");
      } else {
        page.classList.remove("on-screen");
        page.classList.remove("off-left");
        page.classList.add("off-right");
      }
      console.log(currentPageNumber, page, pageNumber);
    });
  }
};

pageElements.nextPage.addEventListener(
  "click",
  pageTurnEventListener(() => currentPageNumber < pages.length - 1, 1),
);

pageElements.previousPage.addEventListener(
  "click",
  pageTurnEventListener(() => currentPageNumber > 0, -1),
);
