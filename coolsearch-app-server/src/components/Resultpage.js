import React from 'react';
import '../css/Resultpage.css';
import FlipMove from 'react-flip-move';
import {Link} from 'react-router-dom';

function Resultpage(props) {
    const items = props.items;
    const resultavailable = props.availableuuid;
    const listItems = items.map(item => {
        const newTo = {
            pathname: "/Result/" + item.key,
            keyword: item.text
        };

        const noLink = {
            pathname: "/",
            keyword: ''
        };

        if(resultavailable === item.key){
            item.status = 'Available'
        }

        if(item.status === 'Available'){
            return <div className="list" key={item.key}>
                <p> {item.text}
                    <span>
                        <Link to={newTo}>Available</Link>
                    </span>
                </p>
            </div>
        } else {
            return (
                <div className="list" key={item.key}>
                    <p> {item.text}
                        <span>
                            <Link to={noLink}>Pending</Link>
                        </span>
                    </p>
                </div>
            )
        }
    })
    return (
      <div className="Recentsearch">
          {listItems}
      </div>
    );
}

export default Resultpage;