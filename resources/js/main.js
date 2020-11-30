const server = "http://localhost:8000/";
const states = getStates();
const data = {
    missed: searchMissedTodoItems() || [],
    valid: searchNotMissedTodoItems() || []
};

/**
 *
 * fonction permettant de récupérer depuis le serveur toutes les valeures possibles de l'état de la tâche
 *
 * @returns {result} : un tableau des valeures
 */
function getStates() {

    let result = null;

    $.ajax({
        type: "GET",
        url: server + "states",
        async: false,
        success: function (data) {
            result = data;
        }
    });

    return result;

}

/**
 *
 * Fonction permettant de supprimer toutes les cartes du DOM sauf celle servant de modèle et celle permettant d'en créer
 *
 */
function resetCards() {

    $('.card').not('.hidden').not('.card-deck > div:last').remove();

}

/**
 * fonction qui parcourt toutes les tâches et les ajoute dans le DOM
 */
function renderTodoList() {

    $.each(data, function (i, e) {
        $.each(e, function (index, elem) {
            addItemTodoList(elem.Id, elem.Nom, elem.Titre, elem.Description, elem.Etat, elem.DateRendu, e !== data.valid);
        })
    })

}

/**
 * Fonction permettant de récupérer toutes les tâches qui ont été manqués et dont le nom correspond au contenu de la barre de recherche
 * ceci est effectué en faisant une requête ajax vers le serveur
 *
 * @returns {result} : un tableau contenant les tâches correspondantes
 */
function searchMissedTodoItems() {

    let result = null;

    $.ajax({
        type: "POST",
        url: server + "searchMissed",
        data: {
            'critere': $("input[type =  search]").val()
        },
        async: false,
        success: function (data) {
            console.log(data)
            result = data;
        }
    });

    return result;

}

/**
 * Fonction permettant de récupérer toutes les tâches encore valides et dont le nom correspond au filtre de la barre de recherche
 * une requête ajax est effectuée vers le serveur pour récupérer ces données
 *
 * @returns {result} : un tableau contenant les tâches correspondantes
 */
function searchNotMissedTodoItems() {

    let result = null;

    $.ajax({
        type: "POST",
        url: server + "searchNotMissed",
        data: {
            'critere': $("input[type =  search]").val()
        },
        async: false,
        success: function (data) {
            console.log(data)
            result = data;
        }
    });

    return result;

}

/**
 *
 * Fonction appelée lors de la modification d'une carte
 * Elle envoie une requête au serveur pour mettre à jour la tâche dans la base de données
 *
 * @param id : int identifiant de la tâche
 * @param newValue : string nouvelle valeure de l'État
 */
function updateTodoItemBackEnd (id, newValue) {

    $.ajax({
        type: "POST",
        url: server + "todo/" + id,
        data: {
            'etat': newValue
        },
        async: false,
        success: function (data) {
            console.log(data)
        }
    });

}

/**
 * Suppression de la tâche du côté serveur
 * Une requête ajax est envoyée au serveur pour supprimer l'item dans la base de donnée
 *
 * @param id : string l'identifiant de la tâche à supprimer
 */
function deleteTodoItemBackEnd(id) {

    $.ajax({
        type: "DELETE",
        url: server + "todo/" + id,
        async: false,
        success: function (data) {
            console.log(data);
        }
    });

}

/**
 * ajoute une carte dans le DOM suivant les paramètres donnés
 * La carte prend alors une couleur suivant son État :
 * \t   -\t à faire : gris
 * \t   -\t en cours : bleu
 * \t   -\t fait : vert
 * \t   -\t rate : noir
 *
 * La date est ensuite converti pour récupérer seulement le jour le mois et l'année et l'afficher correctement en lettre
 *
 * Pour rendre le site un peu plus responsive, en fonction du nombre de carte présentes et de la taille de l'écran, un retour à la ligne est ajouté
 *
 * @param id : int identifiant de la tâche
 * @param nom : string Nom de la tâche
 * @param titre : string titre de la tâche
 * @param contenu : string contenu de la tâche
 * @param state : string État de la tâche
 * @param dateRendu : string date de fin de la tâche
 * @param missed : boolean vrai si la tâche a été manqué
 */
function addItemTodoList (id, nom, titre, contenu, state, dateRendu, missed) {

    const newItem = $('.hidden').first().clone();
    newItem.removeClass("hidden");
    newItem.find('.card-header').text(nom);
    newItem.find('.card-title').text(titre);
    newItem.find('p .card-text').text(contenu);
    newItem.find('select').val(state).change(function () {
        updateTodoItemBackEnd(id, $(this).val());
        updateTodoList();
    });
    newItem.attr('id', id);
    if (missed)
        newItem.find('select').prop('disabled', 'disabled');
    dateRendu = convertDate(dateRendu);
    newItem.removeClass(function (index, className) {
        return (className.match(/(^|\s)bg-\S+/g) || []).join(' ');
    });

    let className;

    switch (state) {
        case states[0]:
            className = "bg-light";
            break;
        case states[1]:
            className = "bg-info"
            break;
        case states[2]:
            className = "bg-success";
            break;
    }

    if (missed)
        className = "bg-dark";

    newItem.addClass(className);

    newItem.find('.card-text > small > em').text(dateRendu);

    const cardDeck = $('.card-deck');
    const nbCard = $('.card').length;

    cardDeck.find(' > div:last').before(newItem);

    if ((nbCard - 1) % 2 === 0) {
        $('.card-deck > div:last').before("<div class=\"w-100 d-none d-sm-block d-md-none\"><!-- wrap every 2 on sm--></div>");
    }
    if ((nbCard - 1) % 3 === 0) {
        $('.card-deck > div:last').before("<div class=\"w-100 d-none d-md-block d-lg-none\"><!-- wrap every 3 on md--></div>");
    }
    if ((nbCard - 1) % 4 === 0) {
        $('.card-deck > div:last').before("<div class=\"w-100 d-none d-lg-block d-xl-none\"><!-- wrap every 4 on lg--></div>");
    }
    if ((nbCard - 1) % 5 === 0) {
        $('.card-deck > div:last').before("<div class=\"w-100 d-none d-xl-block\"><!-- wrap every 5 on xl--></div>");
    }

}

/**
 * Fonction qui converti la date de type string dans un format français en lettre
 *
 * @param date : string La date à convertir
 * @returns {string} : la date converti en lettre
 */
function convertDate(date) {

    if (date.indexOf('T') >= 0) {
        date = date.substring(0, date.indexOf('T'));
    }
    const dateNumber = date.split('-');
    console.log(dateNumber);
    let day = dateNumber[2], year = dateNumber[0], month = dateNumber[1];

    switch (parseInt(month)) {
        case 0:
            month = "Janvier";
            break;
        case 1:
            month = "Février";
            break;
        case 2:
            month = "Mars";
            break;
        case 3:
            month = "Avril";
            break;
        case 4:
            month = "Mai";
            break;
        case 5:
            month = "Juin";
            break;
        case 6:
            month = "Juillet";
            break;
        case 7:
            month = "Août";
            break;
        case 8:
            month = "Septembre";
            break;
        case 9:
            month = "Octobre";
            break;
        case 10:
            month = "Novembre";
            break;
        default:
            month = "Décembre";
    }

    return day + " " + month + " " + year;

}

/**
 *
 * Ajout d'une tâche du côté server à l'aide d'une requête ajax
 *
 * @param nom : string le nom de la tâche
 * @param titre : string le titre de la tâche
 * @param contenu : string le contenu de la tâche
 * @param dateRendu : string la date de fin de la tâche
 */
function addItemToBackEnd (nom, titre, contenu, dateRendu) {

    $.ajax({
        type: "POST",
        url: server + "todo",
        async: false,
        data: {
            'nom': nom,
            'titre': titre,
            'description': contenu,
            'date_rendu': dateRendu
        },
        success: function (data) {
            console.log(data);
        }
    });

}

/**
 *
 * mise à jour du DOM
 *
 */
function updateTodoList() {

    data.missed = searchMissedTodoItems();
    data.valid = searchNotMissedTodoItems();
    resetCards();
    renderTodoList();

}

/**
 * sauvegarde de la nouvelle tâche au clique sur le bouton sauvegarder après une simple vérification du formulaire
 */
$('#sauvegarder').click(function (e) {
    const form = $('#form');
    if (document.getElementById("form").reportValidity()) {
        e.preventDefault();
        addItemToBackEnd($('#nom').val(), $('#titre').val(), $('#contenu').val(), $('#date_rendu').val());
        updateTodoList();
        form.trigger("reset");
    }

});

const $modal = $('#confirmSupression');

/**
 * suppression de la tâche après confirmation du modal
 */
$modal.on("shown.bs.modal", function(event) {
    $modal.on('click', '#supprimer', function () {
        deleteTodoItemBackEnd(event.relatedTarget.parentElement.parentElement.getAttribute('id'));
        $modal.modal("hide");
        updateTodoList();
    });
});

/**
 * filtre des cartes après appuie d'une touche sur le champ de recherche
 */
$('input[type = search]').keyup(function () {
    console.log($(this).val())
    data.missed = searchMissedTodoItems();
    data.valid = searchNotMissedTodoItems();
    updateTodoList();
})

updateTodoList();

setInterval(updateTodoList, 5000);
