import van from "vanjs-core"

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

                    reject()
                })
        
                return;
            }

            resp.json().then(obj => {
                console.log("loaded setting", obj);

                resolve(obj)
            })
        }).catch(err => {
            console.log("failed to make get setting request", err);

            reject()
        })
    })
}

function storeSetting({name, value}) {
    return new Promise((resolve, reject) => {
        const object = {name, value}

        console.log("storing setting", object)

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
    })
}

class Setting {
    constructor({name, defaultValue}) {
        this.name = name
        this.defaultValue = defaultValue
        this.value = null
        this.ready = false

        loadSetting(this.name).then(setting => {
            this.value = setting.value
            this.ready = true
        }).catch(() => {
            this.value = this.defaultValue
            this.ready = true
        })
    }

    setValue(value) {
        this.value = value

        storeSetting(this.name, value)
    }

    isReady() {
        return this.ready
    }

    getValue() {
        return this.value
    }
}

export default Setting