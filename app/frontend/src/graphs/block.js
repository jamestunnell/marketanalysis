import van from "vanjs-core"
import hash from 'object-hash';

import { Button } from '../elements/buttons.js'
import EditBlockForm from './editblock.js'
import { ModalBackground, ModalForeground } from '../modal.js'

class BlockItem {
    constructor({id, config, info, parent}) {
        this.id = id
        this.info = info
        this.parent = parent
        this.type = config.type

        this.name = van.state(config.name)
        this.config = van.state(config)
    }

    getName() {
        return this.name.val
    }
    
    getInfo() {
        return this.info
    }

    makeConfig() {
        return this.config.val
    }

    delete() {
        this.parent.deleteBlock(this.id)
    }

    renderButton() {
        return Button({child: this.name, onclick: () => this.editModal()})
    }

    editModal() {
        const configBefore = this.makeConfig()
        const closed = van.state(false)
        const onComplete = (config) => {
            console.log(`completed block edit`, {config})

            if (hash(config) !== hash(configBefore)) {
                console.log('updating parent graph digest')

                this.config.val = config
                this.name.val = config.name

                this.parent.updateDigest()
            }
    
            closed.val = true
        }
        const onCancel = () => closed.val = true;
        const onDelete  = () => {
            closed.val = true

            this.parent.deleteBlock(this.id)
        }
        const form = EditBlockForm({
            config: configBefore,
            info: this.info,
            otherNames: this.parent.blockNames().filter(name => name !== configBefore.name),
            possibleSources: this.parent.getPossibleSources(configBefore.name),
            onComplete, onCancel, onDelete,
        })
        const modal = ModalBackground(ModalForeground({}, form))
    
        console.log("editing block", configBefore)

        van.add(document.body, () => closed.val ? null : modal);
    }
}

export {BlockItem};