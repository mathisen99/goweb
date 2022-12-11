const likeThumbs = document.getElementsByClassName("js-like")
for (var i = 0; i < likeThumbs.length; i++) {
    const id = likeThumbs[i].dataset.postid;
    likeThumbs[i].addEventListener('click', sendData.bind(null, id, 1), false);
}

const dislikeThumbs = document.getElementsByClassName("js-dislike")
for (var i = 0; i < dislikeThumbs.length; i++) {
    const id = dislikeThumbs[i].dataset.postid;
    dislikeThumbs[i].addEventListener('click', sendData.bind(null, id, 0), false);
}

function sendData(id, isLike) {
    const XHR = new XMLHttpRequest();

    const urlEncodedDataPairs = [];

    // Turn the data object into an array of URL-encoded key/value pairs.
    urlEncodedDataPairs.push(`${encodeURIComponent("postID")}=${encodeURIComponent(id)}`);
    urlEncodedDataPairs.push(`${encodeURIComponent("isLike")}=${encodeURIComponent(isLike)}`);

     // Combine the pairs into a single string and replace all %-encoded spaces to
    // the '+' character; matches the behavior of browser form submissions.
    const urlEncodedData = urlEncodedDataPairs.join('&').replace(/%20/g, '+');

    // Define what happens on successful data submission
    XHR.addEventListener('load', (event) => {
        
        if (isLike) {
            element = document.querySelector(".js-like[data-postid='"+ id +"']")   
        } else {
            element = document.querySelector(".js-dislike[data-postid='"+ id +"']")
        }
        
        if (event.currentTarget.status == 202) {
            if (element.classList.contains("liked")) {
                element.classList.remove("liked");
          
                const countItems = element.closest(".js-like-count").getElementsByClassName("js-count")

                counter = countItems[0].innerHTML 
                countItems[0].innerHTML = parseInt(counter) - 1;
                
            } else {
                element.classList.add("liked");
          
                const countItems = element.closest(".js-like-count").getElementsByClassName("js-count")
        
                counter = countItems[0].innerHTML 
                countItems[0].innerHTML = parseInt(counter) + 1;  
            }
         } else if (event.currentTarget.status == 401) {
                const countItems = element.closest(".js-like-count").getElementsByClassName("js-like-error")
                countItems[0].innerHTML = "<div class=\"like-error\">You need to login to send a reaction!<\/div>";     
        } else {
            alert('Oops! Something went wrong.');
        }
    });

    // Define what happens in case of error
    XHR.addEventListener('error', (event) => {
        alert('Oops! Something went wrong.');
    });

    // Set up our request
    XHR.open('POST', '/add_like');

    // Add the required HTTP header for form data POST requests
    XHR.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');

    // Finally, send our data.
    XHR.send(urlEncodedData);
}


const likeCommentThumbs = document.getElementsByClassName("js-comment-like")
for (var i = 0; i < likeCommentThumbs.length; i++) {
    const id = likeCommentThumbs[i].dataset.commentid;
    likeCommentThumbs[i].addEventListener('click', sendCommentData.bind(null, id, 1), false);
}

const dislikeCommentThumbs = document.getElementsByClassName("js-comment-dislike")
for (var i = 0; i < dislikeCommentThumbs.length; i++) {
    const id = dislikeCommentThumbs[i].dataset.commentid;
    dislikeCommentThumbs[i].addEventListener('click', sendCommentData.bind(null, id, 0), false);
}

function sendCommentData(id, isLike) {
    const XHR = new XMLHttpRequest();

    const urlEncodedDataPairs = [];

    // Turn the data object into an array of URL-encoded key/value pairs.
    urlEncodedDataPairs.push(`${encodeURIComponent("commentID")}=${encodeURIComponent(id)}`);
    urlEncodedDataPairs.push(`${encodeURIComponent("isLike")}=${encodeURIComponent(isLike)}`);

     // Combine the pairs into a single string and replace all %-encoded spaces to
    // the '+' character; matches the behavior of browser form submissions.
    const urlEncodedData = urlEncodedDataPairs.join('&').replace(/%20/g, '+');

    // Define what happens on successful data submission
    XHR.addEventListener('load', (event) => {
        if (isLike) {
            element = document.querySelector(".js-comment-like[data-commentid='"+ id +"']")   
        } else {
            element = document.querySelector(".js-comment-dislike[data-commentid='"+ id +"']")
        }

        if (event.currentTarget.status == 202) {
            if (element.classList.contains("liked")) {
                element.classList.remove("liked");
          
                const countItems = element.closest(".js-like-count").getElementsByClassName("js-count")

                counter = countItems[0].innerHTML 
                countItems[0].innerHTML = parseInt(counter) - 1;
                
            } else {
                element.classList.add("liked");
          
                const countItems = element.closest(".js-like-count").getElementsByClassName("js-count")
        
                counter = countItems[0].innerHTML 
                countItems[0].innerHTML = parseInt(counter) + 1;  
            }

        } else if (event.currentTarget.status == 401) {
            const countItems = element.closest(".js-like-count").getElementsByClassName("js-like-error")
            countItems[0].innerHTML = "<div class=\"like-error\">You need to login to send a reaction!<\/div>";
        } else {
            alert('Oops! Something went wrong.');
        }
    });

    // Define what happens in case of error
    XHR.addEventListener('error', (event) => {
        alert('Oops! Something went wrong.');
    });

    // Set up our request
    XHR.open('POST', '/add_comment_like');

    // Add the required HTTP header for form data POST requests
    XHR.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');

    // Finally, send our data.
    XHR.send(urlEncodedData);
}