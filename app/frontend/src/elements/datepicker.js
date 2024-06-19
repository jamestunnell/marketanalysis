import Datepicker from 'flowbite-datepicker/Datepicker'

import Textbox from './textbox.js'

const WeekdayDatepicker = ({id="", value, maxDate, containerID}) => {
    const input = Textbox({id, value, placeholder: 'Select date'})

    input.addEventListener('changeDate', (e) => value.val = e.target.value)

    const datePickerOpts = {
        autohide: true,
        container: `#${containerID}`,
        daysOfWeekDisabled: [0, 6], // disable saturday and sunday
        format: "yyyy-mm-dd",
        maxDate,
        todayHighlight: true,
    }
    const datepicker = new Datepicker(input, datePickerOpts)

    return input
}

export { WeekdayDatepicker }