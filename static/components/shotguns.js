class ShotgunSpreadCalculator extends HTMLElement {
  constructor() {
    super()
    const template = document.getElementById('shotgun-spread-calculator');

    const templateContent = template.content.cloneNode(true);
    this.appendChild(templateContent);

    this.container = this.querySelector(".hit-locations")
    this.form = this.querySelector(".spread-calc-form")
    this.clearButton = this.querySelector(".clear-locations")
  }

  connectedCallback() {
    this.form.addEventListener('submit', (event) => {
      event.preventDefault();
      const data = new FormData(this.form);

      const poi = data.get("impactPoint");
      const salm = data.get("salm");
      const hc = data.get("hitCount");
      const dt = data.get("diceType")

      this.setLocations(salm, dt, hc, poi)
    });

    this.clearButton.addEventListener('click', (event) => {
      event.preventDefault()
      this.container.replaceChildren()
    })
  }

  clearLocations() {
    this.container.replaceChildren();
  }

  setLocations(salm, dt, hc, poi) {
    this.clearLocations()
    const locations = generateLocations(salm, dt, hc, poi)
    for(let i = 0; i < locations.length; i++) {
      const span = document.createElement('span');
      span.textContent = locations[i];
      this.container.appendChild(span)
    }
  }
}

customElements.define( 'shotgun-spread-calculator', ShotgunSpreadCalculator);

