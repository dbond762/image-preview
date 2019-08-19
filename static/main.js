document.addEventListener('DOMContentLoaded', () => {
  let urlBlockTemplate = `<div class="image-loader__url-block url-block">
                            <input class="url-block__url" type="url">
                            <span class="url-block__add">Add</span>
                            <span class="url-block__remove url-block__remove--disable">Remove</span>
                          </div>`;

  document.querySelector('.image-loader__load').addEventListener('click', () => {
    let inputs = document.querySelectorAll('.url-block__url');

    let urls = [];
    inputs.forEach(url => urls.push(url.value));

    fetch('/preview', {
      method: 'POST',
      cache: 'no-cache',
      headers: {
        'Content-Type': 'applicaton/json',
      },
      body: JSON.stringify({
        urls: urls,
      }),
    })
      .then(response => response.json())
      .then(response => {
        const results = response.previews.map(preview => {
          return `<div class="previews__result preview__wrap">
                    <img class="preview__img" src="${preview}" alt="${preview}">
                  </div>`;
        });

        let previews = document.querySelector('.previews');

        while (previews.firstChild) {
          previews.removeChild(previews.firstChild);
        }

        previews.insertAdjacentHTML('afterBegin', results.join(''));
      })
      .catch(err => {
        console.log(err);
      });
  });
});