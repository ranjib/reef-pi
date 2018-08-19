import React from 'react'
import ColorPicker from './color_picker'
import ProfileSelector from './profile_selector'
import Profile from './profile'

const Channel = (props) => {

  const handleChange = e => {
    props.onChangeHandler(e, props.channelNum)
  }

  const handleConfigChange = e => {
    console.log('handling config change', JSON.stringify(e))
    props.onChangeHandler(e, props.channelNum)
  }

  return (
    <div className="controls border-top">
      <div className="row align-items-start">
        <div className="col-sm-6 col-md-4">
          <div className="form-group">
            <label>Channel Name</label>
            <input type="text" className="form-control" 
              placeholder="Enter Channel Name"
              name="name"
              onChange={handleChange}
              value={props.channel.name} />
          </div>
        </div>
        
        <div className="form-group col-sm-6 col-md-4 form-inline">
          <label className="mb-2">Color</label>
          <ColorPicker name="color"
            color={props.channel.color}
            onChangeHandler={handleChange} />
        </div>

        <div className="col-sm-4">
          <div className="form-group">
            <label>Behavior</label>
            <select className="custom-select"
              name="reverse"
              onChange={handleChange} >
              <option value="false">Active High</option>
              <option value="true">Active Low</option>
            </select>
          </div>
        </div>
      </div>
      <div className="row">
        <div className="col-sm-4">
          <div className="form-group">
            <label>Min</label>
            <input type="text" className="form-control"
              name="min" 
              onChange={handleChange}
              value={props.channel.min} />
          </div>  
        </div>
        <div className="col-sm-4">
          <div className="form-group">
            <label>Max</label>
            <input type="text" className="form-control"
              name="max"
              onChange={handleChange}
              value={props.channel.max}  />
          </div>  
        </div>
        <div className="col-sm-4">
          <div className="form-group">
            <label>Start</label>
            <input type="text" className="form-control"
              name="start_min"
              onChange={handleChange} 
              value={props.channel.start_min} />
          </div>  
        </div>
      </div>
      <div className="row">
        <div className="col">
          <div className="form-group">
            <label className="mr-3">Profile</label>
            <ProfileSelector
              name="profile.type"
              onChangeHandler={handleChange}
              value={props.channel.profile.type} />
          </div>
        </div>
      </div>
      <div className="row mb-3">
        <div className="col">
          <Profile 
            name='profile.config'
            type={props.channel.profile.type}
            config={props.channel.profile.config} 
            onChangeHandler={handleConfigChange} />
        </div>
      </div>
    </div>
  )
}

export default Channel