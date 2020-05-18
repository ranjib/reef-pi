import EditJournal from './edit_journal'
import JournalSchema from './schema'
import { withFormik } from 'formik'

const JournalForm = withFormik({
  displayName: 'JournalForm',
  mapPropsToValues: props => {
    let data = props.data
    if (data === undefined) {
      data = {
        name: '',
        description: '',
        unit: ''
      }
    }
    const values = {
      id: data.id || '',
      name: data.name || '',
      description: data.description || '',
      unit: data.unit || ''
    }
    return values
  },
  validationSchema: JournalSchema,
  handleSubmit: (values, { props }) => {
    props.onSubmit(values)
  }
})(EditJournal)

export default JournalForm
