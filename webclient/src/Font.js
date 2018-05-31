import React, {Component} from 'react';
import './App.css';
import * as request from 'request-promise';
import ReactDynamicFont from 'react-dynamic-font';
import { range } from 'lodash';
import { Motion, spring } from 'react-motion';


const config = {
    server: {
        method: 'http',
        hostname: 'localhost',
        port: 8080,
        get baseUrl() {
            return `${config.server.method}://${config.server.hostname}:${config.server.port}`;
        }
    }
};
export class Font extends Component {
    constructor(props) {
        super(props);
        this.state = {
            text: " ",
            board: null,
        }

        this.handleKeyPress = this.handleKeyPress.bind(this)
    }

    async componentDidMount() {
        const board = await request.post({
            url: `${config.server.baseUrl}/v1/fonts/render`,
            json: true,
            body: {
                "fontName": "TI84",
                "text": "",
                "spaceWidth": 0,
                "kerning": 0
            }
        });

        this.setState({
            board: board,
        });
        console.log("setting board to: ", board)
    }

    handleKeyPress(event) {
        this.setState({
            text: event.target.value
        });

        return request.post({
            url: `${config.server.baseUrl}/v1/fonts/render`,
            json: true,
            body: {
                "fontName": "TI84",
                "text": event.target.value,
                "spaceWidth": 1,
                "kerning": 0
            }
        })
            .then((board) => {
                this.setState({
                    board: board
                })
            })
        console.log(2)
        console.log(3)
    }

    render() {
        //Concept behind this was to create a board that replicates the physical board
        //User can input a text response and the server will fire back by flipping each dot to their corresponding letter

        //this will return an empty div and generate a blank board

        // const emptyBoard = [];
        //
        //
        //
        // const board = (numberOfRows, numberOfColumns) => {
        //
        //     return (<div>
        //         {(range(numberOfRows)).map((blankDot, x) => {
        //             if (emptyBoard.length < numberOfRows) {
        //                 emptyBoard.push([]);
        //             }
        //             return (<div>
        //                 {range(numberOfColumns).map((dot, y) => {
        //                     if (emptyBoard.length < numberOfColumns) {
        //                         // emptyBoard[x].push(Math.floor(Math.random() * 5) % 2 === 0);
        //                     }
        //                     console.log(emptyBoard[x][y] ? "frontSide" : "backSide")
        //                     // return (<span className={`disc ${emptyBoard[x][y] ? "frontSide" : "backSide"}`}></span>);
        //                     return (<span className={`disc ${
        //                         drawABoard("⚪️") ? "frontSide": "backSide"}`}></span>);
        //                 })}
        //             </div>)
        //         })}
        //     </div>)
        // }


        const drawABoard = () => {
            if (!this.state.board) {
                return (<div></div>)
            }

            // for each row, create a span of columns whiteSpace={"wrap"}
            return (

                <div className="board-container">
                    <ReactDynamicFont smooth content={ this.state.board.map((line, lineNumber) => {
                        return (<div className="line-container">
                            {line.map((letter, letterIndex) => {
                                return (
                                    <span className="letters">
                                        {
                                            letter.map((boardRow, boardIndex) => {
                                                return (
                                                    <div className="board-row" key={boardIndex}>

                                                        {boardRow.map(dot => <span className=""> <div className="dot-container">{dot}</div></span>)}

                                                    </div>
                                                );
                                            })
                                        }
                                        </span>
                                )
                            })}
                        </div>)
                    })} />

                </div>
            );
        };

        if (!this.state) {
            return (<div></div>)
        }
        return (
            <div>

            <div className="font-container">
                <div>
                    <h4> Try it out here! Type something in the box below! </h4>
                </div>
                <div>
                    <textarea type="text"
                        onChange={this.handleKeyPress}
                        value={this.state.text}>
                    </textarea>
                    {drawABoard()}
                </div>

            </div>
                {/*<div className="disc-container">*/}
                    {/*EMPTY BOARD WILL GO HERE:*/}

                    {/*{board(14, 28)}*/}


                {/*</div>*/}
            </div>

        )
    }
}

export default Font;