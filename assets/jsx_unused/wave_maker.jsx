import React from 'react';

export default class WaveMaker extends React.Component {
    constructor(props) {
        super(props);
    }
    onChange(){
    }

    render() {
      var borderStyle = {
        borderWidth: '1px',
        borderStyle: 'solid',
        borderColor: '#dddddd'
      };
      return (
          <div className="container">
            <div className="row">
              Run after <input type="text" value="4h" onChange={this.onChange}/> for <input type="text" value="4m" onChange={this.onChange}/>
            </div>
            <div className="row">
              <div className="col-sm-4">Setup wave maker</div>
              <div className="col-sm-4"><input type="button" value="On/Off"/></div>
            </div>
          </div>
          );
    }
}
