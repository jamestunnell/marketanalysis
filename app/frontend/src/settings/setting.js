import van from 'vanjs-core'
import storeSetting from './storesetting'

class Setting {
    constructor({name, storedValue, defaultValue}) {
        this.name = name
        this.storedValue = storedValue
        this.value = van.state(storedValue ?? defaultValue)

        van.derive(() => {
            if (this.storedValue && this.value.val === this.storedValue) {
                return
            }

            const s = {name: this.name, value: this.value.val}
            
            storeSetting(s).then(() => {
                this.storedValue = this.value.val
            })
        })
    }
}

export default Setting