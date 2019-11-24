import React from 'react';
import ReactDOM from 'react-dom';
import '../css/Searchresult.css';

const kafkaHost = process.env.SERVICE_HOST
if ('SERVICE_HOST' not in process.env){
  kafkaHost = "localhost:9092"
}

function componentDidMount(uuid, keyword) {
const searchWord = <div><center><h1 className="card-title1">Search results for: {keyword}</h1></center></div>
    
    fetch('http://' + kafkaHost + '/v1/api/searchresult/' + uuid)
    .then(res => res.json())
    .then((data) => { 
       var googleresult = data["google"].map((google) => (
            <div className="card">
                <div className="card-body">
                    <a target="_blank" href={google.link}><h6 className="card-title">{google.title}</h6></a>
                    <p className="card-subtitle mb-2 text-muted">{google.snippet}</p>
                </div>
            </div>
        ))
        var ddgeresult = data["ddg"].map((ddg) => (
            <div className="card">
                <div className="card-body">
                    <a target="_blank" href={ddg.link}><h6 className="card-title">{ddg.title}</h6></a>
                    <p className="card-subtitle mb-2 text-muted">{ddg.snippet}</p>
                </div>
            </div>
        ))
        var wikieresult = data["wiki"].map((wiki) => (
            <div className="wikicard">
                <div className="card-body">
                    <a target="_blank" href={wiki.link}><h6 className="card-title">{wiki.title}</h6></a>
                    <p className="card-subtitle mb-2 text-muted">{wiki.snippet}</p>
                </div>
            </div>
        ))
        ReactDOM.render(searchWord,document.getElementById("search-word"))
        ReactDOM.render(wikieresult,document.getElementById("wiki-results"))
        ReactDOM.render(googleresult,document.getElementById("google-results"))
        ReactDOM.render(ddgeresult,document.getElementById("ddg-results"))
    })
    .catch(console.log)
}
function Searchresult(props){
    var uuid = props.match.params.uuid;
    var keyword = props.location.keyword;
    console.log("UUID=", uuid, "KeyWord=", keyword)
    componentDidMount(uuid, keyword)
    return(
        <div className="searchresult">
            <div id="search-word" className="searchword"></div>
            <div id="wiki-results" className="card wikicard"></div>
            <section>
                <div className="googleresult">
                    <h5 className="card-title1">Google</h5>
                    <div id="google-results"/>
                </div>
                <div className="ddgresults">
                    <h5 className="card-title1">Duck Duck Go</h5>
                    <div id="ddg-results"/>
                </div>
            </section>
        </div>
    )
}

export default Searchresult;