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

  submitButton: document.querySelector("form input[type=submit]"),
  thanksForSubmitting: document.querySelector(
    "#appointment-form-submitted > h2",
  ),
};

const appointmentRequestForm = document.querySelector(
  "#appointment-request-form",
);

const appointmentFormSubmitted = document.querySelector(
  "#appointment-form-submitted",
);

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
  submitButton: "Submit",
  thanksForSubmitting:
    "Thank you for requesting your appointment to analyze your home water supply! We will contact you as soon as we are able.",
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
  submitButton: "Enviar",
  thanksForSubmitting:
    "¡Gracias por solicitar su cita para analizar el agua en su hogar! Nos pondremos en contacto con usted lo antes posible.",
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
    if (elementName === "submitButton") {
      // the submit button must be handled as a special case
      // because its text is in its "value" attribute
      // and not in its textContent
      element.setAttribute("value", text[elementName]);
    } else {
      element.textContent = text[elementName];
    }
  });
});

let currentPageNumber = 0;

const getPages = () => {
  // console.log("getting pages now 🔍");
  const pages = [...document.querySelectorAll("[class*=page-name]")];
  // console.log(pages);
  return pages;
};

const numberOfBackgroundColors = 4;

const pageTurnEventListener = (predicate, pageIncrement) => () => {
  if (predicate()) {
    const pages = getPages();
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

    document.body.classList.remove(
      `color-${(currentPageNumber - pageIncrement) % numberOfBackgroundColors}`,
    );
    document.body.classList.add(
      `color-${currentPageNumber % numberOfBackgroundColors}`,
    );
    pages.forEach((page, pageNumber) => {
      // TODO: this is a little ugly 😬
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
  pageTurnEventListener(() => currentPageNumber < getPages().length - 1, 1),
);

pageElements.previousPage.addEventListener(
  "click",
  pageTurnEventListener(() => currentPageNumber > 0, -1),
);

// page turn with arrow keys left and right
document.body.onkeydown = ({ code }) => {
  if (code === "ArrowLeft") pageElements.previousPage.click();
  if (code === "ArrowRight") pageElements.nextPage.click();
};

// page turn with touch swipe
const touches = {};
document.addEventListener(
  "touchstart",
  ({ changedTouches: [{ clientX, clientY, identifier }] }) =>
    (touches[identifier] = [clientX, clientY]),
);

document.addEventListener(
  "touchend",
  ({ changedTouches: [{ clientX, clientY, identifier }] }) => {
    const [startX, startY] = touches[identifier];
    const deltaX = clientX - startX;
    const deltaY = clientY - startY;
    if (
      Math.abs(deltaX) > Math.abs(deltaY) &&
      Math.abs(deltaX) > window.innerWidth / 3
    ) {
      if (deltaX < 0)
        // swiped right
        pageElements.nextPage.click();
      if (deltaX > 0)
        // swiped left
        pageElements.previousPage.click();
    }
  },
);

[...document.querySelectorAll("form input")].forEach((inputElement) => {
  inputElement.onkeydown = (e) => e.stopPropagation();
});

appointmentRequestForm.addEventListener("submit", (submissionEvent) => {
  submissionEvent.preventDefault();
  console.log("submitted");
  console.log(submissionEvent);
  pageElements.submitButton.disabled = true;
  const formFields = [...submissionEvent.target].filter(
    ({ tagName }) => tagName !== "FIELDSET",
  );
  fetch("https://hqlpruvski.execute-api.us-west-2.amazonaws.com/appt", {
    method: "POST",
    body: JSON.stringify({
      name: formFields[0].value,
      addressLine1: formFields[1].value,
      addressLine2: formFields[2].value,
      city: formFields[3].value,
      state: formFields[4].value,
      zip: formFields[5].value,
      phone: formFields[6].value,
      apptDateAndTime: formFields[7].value,
    }),
  })
    .then((r) => {
      if (r.status === 200) {
        appointmentFormSubmitted.classList = appointmentRequestForm.classList;
        appointmentRequestForm.classList = "inactive";
        localStorage.setItem("has-submitted", true);
      } else {
        console.error(r);
      }
    })
    .catch((e) => {
      console.error(e);
    });
});

if (localStorage.getItem("has-submitted")) {
  appointmentFormSubmitted.classList = appointmentRequestForm.classList;
  appointmentRequestForm.classList = "inactive";
}

// section for the first animation 😁

const polarToPolygon = (theta, r, xoffset = 0, yoffset = 0) => {
  const x = Math.cos(theta) * r * 50;
  const y = Math.sin(theta) * r * 50;
  return [`${x + 50 + xoffset}%`, `${y + 50 + yoffset}%`];
};

const spikyClipPath = (spikeCount, innerRadius, spikeLength) =>
  `polygon(${new Array(2 * spikeCount)
    .fill()
    .map((_, i) => [
      ((1 + i) * 2 * Math.PI) / (2 * spikeCount),
      innerRadius + spikeLength * (i % 2),
    ])
    .map(([a, b]) => polarToPolygon(a, b))
    .flatMap((e) => e.join(" "))
    .join(", ")})`;

[...document.querySelectorAll("div[class*=animation001] > div")].forEach(
  (e, i) => {
    const spiky = document.createElement("div");
    const spikeCount = Math.floor(36 * Math.random() + 4);
    const innerRadius = Math.random();
    const spikeLength = 1 - innerRadius;
    spiky.style.setProperty(
      "clip-path",
      spikyClipPath(spikeCount, innerRadius, spikeLength),
    );
    const hue = Math.floor(360 * Math.random());
    spiky.style.setProperty(
      "background-image",
      `linear-gradient(45deg, hsl(${hue}deg, 100%, 50%), hsl(${hue + 90}deg, 100%, 50%))`,
    );
    const animationDuration = Math.random() * 10 + 5;
    spiky.style.setProperty("animation-duration", `${animationDuration}s`);
    spiky.style.setProperty("animation-timing-function", "linear");
    spiky.style.setProperty("animation-name", "spinning");
    spiky.style.setProperty("animation-iteration-count", "infinite");
    spiky.style.setProperty("position", "absolute");
    spiky.style.setProperty("top", "-0.5em");
    spiky.style.setProperty("left", "0px");
    spiky.style.setProperty("width", "2.5em");
    spiky.style.setProperty("height", "2.5em");
    spiky.style.setProperty("z-index", `-${i + 1}`);
    e.style.setProperty("filter", `drop-shadow(20px 20px 4px hsl(${hue + 180}deg, 100%, 50%))`);
    const fontSize = Math.random() * 4 + 4;
    e.style.setProperty("font-size", `${fontSize}rem`);
    e.style.setProperty("z-index", "1");
    e.style.setProperty("position", "relative");
    e.style.setProperty("left", `${4 * (i % 2)}em`);
    e.appendChild(spiky);
  },
);
