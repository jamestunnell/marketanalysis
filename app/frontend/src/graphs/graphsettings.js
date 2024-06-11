import van from "vanjs-core"
import Datepicker from 'flowbite-datepicker/Datepicker'

import { Get, PostJSON, PutJSON } from "../backend"

const { input } = van.tags

const INPUT_CLASS = "px-1 py-1 text-gray-500 rounded-md focus:ring-gray-500 focus:border-gray-500 focus:ring"

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
        this.numChartsInput = input({
            id: "numCharts",
            class: INPUT_CLASS,
            type: "number",
            min: 1,
            step: 1,
            onchange: e => {
                const strVal = e.target.value
                const newVal = parseInt(strVal)
                if (isNaN(newVal)) {
                    console.log(`ignoring NaN numCharts ${strVal}`)

                    return
                }

                this.numCharts.val = newVal

                storeSetting({name: "numCharts", value: newVal})
            },
        })
        this.symbolInput = input({
            id: "symbol",
            class: INPUT_CLASS,
            type: "text",
            placeholder: "SPY, QQQ, etc.",
            onchange: e => {
                this.symbol.val = e.target.value

                storeSetting({name: "symbol", value: e.target.value})
            },
        })
        this.dateInput = input({
            id: "actionDate",
            class: INPUT_CLASS,
            type: "text",
            placeholder: 'Select date',
        });

        this.dateInput.addEventListener('changeDate', (e) => {
            console.log(`changed date to ${e.target.value}`, e)

            storeSetting({name: "date", value: e.target.value}).then(() => {
                this.date.val = e.target.value
            })
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