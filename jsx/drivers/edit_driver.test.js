import React from 'react'
import * as Alert from '../utils/alert'
import Enzyme, {shallow } from 'enzyme'
import Adapter from 'enzyme-adapter-react-16'
import EditDriver from './edit_driver'
import 'isomorphic-fetch'

Enzyme.configure({ adapter: new Adapter() })

describe('<EditDriver />', () => {
  var values = {}

  var fn = jest.fn()

  beforeEach(() => {
    jest.spyOn(Alert, 'showError')
  })

  afterEach(() => {
    jest.clearAllMocks()
  })

  it('<EditDriver />', () => {
    shallow(
      <EditDriver values={values}
        handleBlur={fn}
        submitForm={fn}
      />
    )
  })
})
