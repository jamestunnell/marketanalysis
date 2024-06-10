import van from "vanjs-core"
import hash from 'object-hash';

import { Button } from '../buttons.js';
import { EditBlockForm } from './editblock.js'
import { ModalBackground } from "../modal.js";

const {div} = van.tags

class BlockItem {
    constructor({id, block, info, parent}) {
        this.id = id
        this.info = info
        this.parent = parent
        this.type = block.type

        this.block = van.state(block)
        this.name = van.derive(() => this.block.val.name)
    }

    getName() {
        return this.name.val
    }
    
    makeBlock() {
        return this.block.val
    }

    delete() {
        this.parent.deleteBlock(this.id)
    }

    renderButton() {
        return Button({child: this.name, onclick: () => this.editModal()})
    }

    editModal() {
        const blockBefore = this.block.val
        const closed = van.state(false)
        const onComplete = (block) => {
            if (hash(block) !== hash(blockBefore)) {
                this.block.val = block
                
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
            block: blockBefore,
            info: this.info,
            otherNames: this.parent.blockNames().filter(name => name !== blockBefore.name),
            possibleSources: this.parent.getPossibleSources(),
            onComplete, onCancel, onDelete,
        })
        const modal = ModalBackground(
            div(
                {class: "block p-8 rounded-lg bg-white z-11"},
                form,
            ),
        )
    
        console.log("editing block", blockBefore)

        van.add(document.body, () => closed.val ? null : modal);
    }
}

export {BlockItem};