class None {
    constructor() {}

    toString() {
        return ""
    }

    validate(val) {
        return null
    }

    isRange() {return true}
    getMin() {return null}
    getMax() {return null}
}

class OneOf {
    constructor(allowed) {
        this.allowed = allowed
    }

    toString() {
        return `${this.allowed}`
    }

    validate(val) {
        if (limits.indexOf(val) === -1) {
            return new Error(`${val} is not one of ${this.allowed}`) 
        }

        return null
    }

    isRange() {return false}
    getMin() {return null}
    getMax() {return null}
}

class Less {
    constructor(max) {
        this.max = max
    }

    toString() {
        return `< ${this.max}`
    }

    validate(val) {
        return (val < this.max) ? null : new Error(`${val} is not < ${this.max}`)
    }

    isRange() {return true}
    getMin() {return null}
    getMax() {return {value: this.max, inclusive: true}}
}

class LessEqual {
    constructor(max) {
        this.max = max
    }

    toString() {
        return `<= ${this.max}`
    }

    validate(val) {
        return (val <= this.max) ? null : new Error(`${val} is not <= ${this.max}`)
    }

    isRange() {return true}
    getMin() {return null}
    getMax() {return {value: this.max, inclusive: true}}
}

class Greater {
    constructor(min) {
        this.min = min
    }

    toString() {
        return `> ${this.min}`
    }

    validate(val) {
        return (val > this.min) ? null : new Error(`${val} is not > ${this.min}`)
    }

    isRange() {return true}
    getMin() {return {value: this.min, inclusive: false}}
    getMax() {return null}
}

class GreaterEqual {
    constructor(min) {
        this.min = min
    }

    toString() {
        return `>= ${this.min}`
    }

    validate(val) {
        return (val >= this.min) ? null : new Error(`${val} is not >= ${this.min}`)
    }

    isRange() {return true}
    getMin() {return {value: this.min, inclusive: true}}
    getMax() {return null}
}

class RangeIncl{
    constructor(min, max) {
        this.min = min
        this.max = max
    }

    toString() {

    }

    validate(val) {
        if (val >= this.min && val <= this.max) {
            return null
        }

        return new Error(`${val} is not in range [${this.min}, ${this.max}]`) 
    }

    isRange() {return true}
    getMin() {return {value: this.min, inclusive: true}}
    getMax() {return {value: this.max, inclusive: true}}
}

class RangeExcl{
    constructor(min, max) {
        this.min = min
        this.max = max
    }
    
    toString() {
        return `[${this.min}, ${this.max})`
    }

    validate(val) {
        if (val >= this.min && val < this.max) {
            return null
        }

        return new Error(`${val} is not in range [${this.min}, ${this.max})`) 
    }

    isRange() {return true}
    getMin() {return {value: this.min, inclusive: true}}
    getMax() {return {value: this.max, inclusive: false}}
}

const MakeConstraint = ({type, limits}) => {
    switch (type) {
    case 'oneOf':
        return new OneOf(limits)
    case 'less':
        return new Less(limits[0])
    case 'lessEqual':
        return new LessEqual(limits[0])
    case 'greater':
        return new Greater(limits[0])
    case 'greaterEqual':
        return new GreaterEqual(limits[0])
    case 'rangeIncl':
        return new RangeIncl(limits[0], limits[1])
    case 'rangeExcl':
        return new RangeExcl(limits[0], limits[1])
    }

    return new None()
}

export {None, Less, LessEqual, Greater, GreaterEqual, RangeExcl, RangeIncl, MakeConstraint}