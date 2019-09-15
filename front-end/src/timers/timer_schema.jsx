import * as Yup from 'yup'
import i18next from 'i18next'

const TimerSchema = Yup.object().shape({
  name: Yup.string()
    .required(i18next.t('timers:name_required')),
  enable: Yup.bool()
    .required(i18next.t('timers:timer_status_required')),
  type: Yup.string()
    .required(i18next.t('timers:type_required')),
  day: Yup.string()
    .required(i18next.t('timers:day_required')),
  hour: Yup.string()
    .required(i18next.t('timers:hour_required')),
  minute: Yup.string()
    .required(i18next.t('timers:minute_required')),
  second: Yup.string()
    .required(i18next.t('timers:second_required')),
  target: Yup.object().shape({
   id: Yup.string(),
   on: Yup.bool(),
   duration: Yup.number(),
   revert: Yup.bool(),
   title: Yup.string(),
   message: Yup.string(),
  }),
})

export default TimerSchema
