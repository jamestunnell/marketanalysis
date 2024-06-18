import van from "vanjs-core"
import Datepicker from 'flowbite-datepicker/Datepicker'

import { Get, PostJSON, PutJSON } from '../backend.js'
import { IntRange } from '../elements/number.js'
import { GreaterEqual } from '../constraint.js'
import Textbox from '../elements/textbox.js'

function loadSetting(name) {
    return new Promise((resolve, reject) => {
        console.log(`loading setting ${name}`)

        Get(`/settings/${name}`).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to load setting", appErr);

                    reject(appErr)
                })
        
                return;
            }

            resp.json().then(obj => {
                console.log("loaded setting", obj);

                resolve(obj)
            })
        }).catch(err => {
            reject({
                title: "Action Failed",
                message: "failed to make get request",
                details: [err.message],
            })
        })
    })
}


function storeSetting({name, value}) {
    return new Promise((resolve, reject) => {
        const object = {name, value}

        console.log("storing setting", object)

        PostJSON({route: "/settings", object}).then(postResp => {
            if (postResp.status === 204) {
                // Avoid Fetch failed loading
                postResp.text().then(text => {
                    resolve()
                })
            } else {
                PutJSON({route: `/settings/${name}`, object}).then(putResp => {
                    if (putResp.status != 204) {
                        putResp.json().then(appErr => {
                            console.log("failed to store setting", appErr);
    
                            reject(appErr)
                        })
                
                        return;
                    }

                    console.log("stored setting", object)
    
                    // Avoid Fetch failed loading
                    putResp.text().then(text => {
                        resolve()
                    })
                }).catch(err => {
                    reject({
                        title: "Action Failed",
                        message: `failed to make put request`,
                        details: [err.message],
                    })
                })
            }
        }).catch(err => {
            reject({
                title: "Action Failed",
                message: `failed to make post request`,
                details: [err.message],
            })
        })
    })
}

class GraphSettings {
    constructor({containerID}) {
        this.date = van.state('')
        this.symbol = van.state('')
        this.numCharts = van.state(1)
        this.containerID = containerID
        this.numChartsInput = IntRange({
            id: "numCharts",
            constraint: new GreaterEqual(1),
            value: this.numCharts,
        })
        this.symbolInput = Textbox({
            id: "symbol",
            placeholder: "SPY, QQQ, etc.",
            value: this.symbol,
        })
        this.dateInput = Textbox({
            id: "actionDate",
            placeholder: 'Select date',
            value: this.date,
        })
        
        this.dateInput.addEventListener('changeDate', (e) => this.date.val = e.target.value)

        van.derive(() => {
            console.log(`changed numCharts to ${this.numCharts.val}`)

            storeSetting({name: "numCharts", value: this.numCharts.val})
        })
        van.derive(() => {
            console.log(`changed symbol to ${this.symbol.val}`)
            
            storeSetting({name: "symbol", value: this.symbol.val})
        })
        van.derive(() => {
            console.log(`changed date to ${this.date.val}`)

            storeSetting({name: "date", value: this.date.val})
        })

        const datePickerOpts = {
            autohide: true,
            container: this.containerID,
            daysOfWeekDisabled: [0, 6], // disable saturday and sunday
            format: "yyyy-mm-dd",
            maxDate: new Date(), // today
            todayHighlight: true,
        }
        const datepicker = new Datepicker(this.dateInput, datePickerOpts)
    }

    load() {
        loadSetting("symbol").then(setting => {
            this.symbolInput.value = setting.value
            this.symbol.val = setting.value
        })

        loadSetting("date").then(setting => {
            this.dateInput.value = setting.value
            this.date.val = setting.value
        })

        loadSetting("numCharts").then(setting => {
            this.numChartsInput.value = setting.value
            this.numCharts.val = setting.value
        })
    }
}

export default GraphSettings