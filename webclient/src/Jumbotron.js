import React,{Component} from 'react';
import './App.css'
import {Jumbotron} from 'react-bootstrap';
import Typist from 'react-typist';
import Font from './Font';
import GithubCorner from 'react-github-corner';



class Jumbo extends Component {
    render() {
        return (
            <div>
            <GithubCorner  href="https://github.com/armory/flipdisks" />
            <Jumbotron className="jumbotronStyle">
                <Typist>
                    <h1>Hello, world!</h1>
                    <p>
                        This is a simple flip disc simulator. =D
                    </p>
                    <p>
                        {/*<Button href="#font" bsStyle="primary">Try out our demo!</Button>*/}
                    </p>
                </Typist>
                <Font/>
            </Jumbotron>
            </div>
        )
    }
}


export default Jumbo;