/*
var dag = {
  nodes: [],
  edges: [],
  questionnaire: {}
  given: []
};
*/

function updateGraph() {
  const graphContainer = d3.select('#graphContainer');

  // Clear the container
  graphContainer.html('');
  graphContainer.append('h2').text('Probability graph');
  const nodeTable = graphContainer.append('table').attr('class', 'table');
  // Create the table header
  const tableHeader = nodeTable.append('thead').append('tr');
  tableHeader.append('th').text('Parents');
  tableHeader.append('th').text('Questioned');
  tableHeader.append('th').text('Children');
  // Create the table body
  const tableBody = nodeTable.append('tbody');
  dag.nodes.forEach(
    (node, node_index) => {
      const tableRow = tableBody.append('tr');
      // Parent Node column
      const parentNodeCell = tableRow.append('td').style('padding', '10px').style('vertical-align', 'middle');
      const parents = [];
      const children = [];
      dag.edges.forEach(
        (edge, index) => {
          if (edge.dst === dag.nodes[node_index].value) {
            parents.push(edge.src);
          }
          if (edge.src === dag.nodes[node_index].value) {
            children.push(edge.dst);
          }
        },
      );
      if (parents && parents.length > 0) {
        parents.forEach(
          (parent, parentIndex) => {
            parentNodeCell.append('span').text(parent);
            parentNodeCell.append('br');
          },
        );
      } // Add the "+" button for selecting the parent node
      const addButton = parentNodeCell.append('button').attr('class', 'btn btn-primary btn-sm ml-2').attr('style', 'margin-left: 0 !important').style('width', '20%').text('+');
      addButton.on(
        'click',
        () => {
          // Create a modal pop-up
          const modal = d3.select('body').append('div').attr('class', 'modal').style('display', 'block')
            .style('position', 'fixed')
            .style('z-index', '9999')
            .style('left', '0')
            .style('top', '0')
            .style('width', '100%')
            .style('height', '100%')
            .style('overflow', 'auto')
            .style('background-color', 'rgba(0, 0, 0, 0.4)');
          // Create the modal content
          const modalContent = modal.append('div').attr('class', 'modal-content').style('background-color', '#fefefe').style('margin', '15% auto')
            .style('padding', '20px')
            .style('border', '1px solid #888')
            .style('width', '50%');
          // Add a close button to the modal
          modalContent.append('span').attr('class', 'close').style('float', 'right').style('font-size', '28px')
            .style('font-weight', 'bold')
            .style('cursor', 'pointer')
            .html('&times;')
            .on(
              'click',
              () => {
                modal.style('display', 'none');
                modal.remove();
              },
            );
          // Add a select dropdown to the modal content
          const selectDropdown = modalContent.append('select').attr('class', 'form-control').style('margin-bottom', '20px');
          // Add options to the select dropdown
          const parentOptions = dag.nodes.filter(
            (node) =>
              // Filter out the current node and its children as options for the parent node
              node.value !== dag.nodes[node_index].value
              && !children.includes(node.value)
            ,
          );
          selectDropdown.selectAll('option').data(parentOptions).enter().append('option')
            .attr('value', (d) => d.value)
            .text((d) => d.value);
          // Add a button to add the edge when selected
          modalContent.append('button').attr('class', 'btn btn-primary').text('Add Edge').on(
            'click',
            () => {
              const selectedParentNode = selectDropdown.node().value;
              if (selectedParentNode) {
                const edge = {
                  src: selectedParentNode,
                  dst: dag.nodes[node_index].value,
                };
                dag.edges.push(edge);
                updateGraph();
                modal.style('display', 'none');
                modal.remove();
              }
            },
          );
        },
      );
      // Node column
      const nodeCell = tableRow.append('td').style('padding', '10px').style('vertical-align', 'middle');
      const deleteButton = nodeCell.append('button').attr('class', 'btn btn-danger btn-sm ml-2').text('-');
      deleteButton.on('click', () => {
        deleteNode(node_index);
      });
      nodeCell.append("div").style('display', 'inline').style('margin-left', '10px').text(node.value);
      // Child Node column
      const childNodeCell = tableRow.append('td').style('padding', '10px').style('vertical-align', 'middle');
      if (children && children.length > 0) {
        children.forEach(
          (child, childIndex) => {
            var childSpan = childNodeCell.append('span').style('margin-bottom', '2px');
            // Create delete button for each child
            const deleteButton = childSpan.append('button').attr('class', 'btn btn-danger btn-sm ml-2').text('-');
            deleteButton.on('click', () => {
              deleteEdge(node_index, child);
            });
            childSpan.append("div").style('display', 'inline').style('margin-left', '10px').text(child);
            if (childIndex < children.length - 1) {
              childNodeCell.append('br');
            }
          },
        );
      } else {
        childNodeCell.text('-');
      }
    },
  );

  // Add the new section for setting Nodes as "given"
  const givenSection = graphContainer.append('div').attr('class', 'section');
  givenSection.append('h3').text('Given Events');

  // Create the table body for "given" nodes
  const givenTable = givenSection.append('table').attr('class', 'table');
  const givenTableBody = givenTable.append('tbody');

  // Add rows for each "given" node
  dag.given.forEach((givenNode) => {
    const givenTableRow = givenTableBody.append('tr');

    // Node cell
    givenTableRow.append('td').text(givenNode);

    // Remove button cell
    const removeButtonCell = givenTableRow.append('td');
    removeButtonCell.append('button')
      .attr('class', 'btn btn-danger btn-sm ml-2')
      .text('Remove')
      .on('click', () => {
        removeGivenNode(givenNode);
      });
  });

  const addGivenButton = givenTableBody.append('tr').append('td').append('button')
    .attr('class', 'btn btn-primary')
    .text('Add Given Node')
    .on('click', addGivenNode);

  graphContainer.append('button').attr('class', 'btn btn-primary')
    .attr('style', 'margin-top: 7px')
    .text('Generate questionnaire').attr('onclick', 'postDag()');
}
function deleteNode(nodeIndex) {
  // Remove the node from the "given" nodes list
  const givenIndex = dag.given.indexOf(dag.nodes[nodeIndex].value);
  if (givenIndex !== -1) {
    dag.given.splice(givenIndex, 1);
  }

  // Remove the edges connected to the node
  dag.edges = dag.edges.filter(
    (edge) => edge.src !== dag.nodes[nodeIndex].value
      && edge.dst !== dag.nodes[nodeIndex].value,
  );

  // Remove the node from the nodes array
  dag.nodes.splice(nodeIndex, 1);

  updateGraph();
}
function deleteEdge(nodeIndex, childNode) {
  // Find the edge to delete
  const edgeIndex = dag.edges.findIndex(
    (edge) => edge.src === dag.nodes[nodeIndex].value
      && edge.dst === childNode,
  );
  if (edgeIndex !== -1) {
    // Remove the edge from the edges array
    dag.edges.splice(edgeIndex, 1);
    // Update the graph
    updateGraph();
  }
}

function addGivenNode() {
  // Create a modal pop-up
  const modal = d3.select('body').append('div').attr('class', 'modal').style('display', 'block')
    .style('position', 'fixed')
    .style('z-index', '9999')
    .style('left', '0')
    .style('top', '0')
    .style('width', '100%')
    .style('height', '100%')
    .style('overflow', 'auto')
    .style('background-color', 'rgba(0, 0, 0, 0.4)');

  // Create the modal content
  const modalContent = modal.append('div').attr('class', 'modal-content').style('background-color', '#fefefe').style('margin', '15% auto')
    .style('padding', '20px')
    .style('border', '1px solid #888')
    .style('width', '50%');

  // Add a close button to the modal
  modalContent.append('span').attr('class', 'close').style('float', 'right').style('font-size', '28px')
    .style('font-weight', 'bold')
    .style('cursor', 'pointer')
    .html('&times;')
    .on(
      'click',
      () => {
        modal.style('display', 'none');
        modal.remove();
      },
    );

  // Add a select dropdown to the modal content
  const selectDropdown = modalContent.append('select').attr('class', 'form-control').style('margin-bottom', '20px');

  // Filter out the nodes that are already given and add them as options to the select dropdown
  const availableNodes = dag.nodes.filter((node) => !dag.given.includes(node.value));

  // Add options to the select dropdown
  selectDropdown.selectAll('option').data(availableNodes).enter().append('option')
    .attr('value', (d) => d.value)
    .text((d) => d.value);

  // Add a button to add the selected node as "given"
  modalContent.append('button')
    .attr('class', 'btn btn-primary')
    .text('Add Node as Given')
    .on('click', () => {
      const selectedNode = selectDropdown.node().value;
      if (selectedNode) {
        dag.given.push(selectedNode);
        updateGraph();
        modal.style('display', 'none');
        modal.remove();
      }
    });
}

function removeGivenNode(node) {
  const index = dag.given.indexOf(node);
  if (index !== -1) {
    dag.given.splice(index, 1);
    updateGraph();
  }
}

function postDag() {
  const data = JSON.stringify({
    nodes: dag.nodes,
    edges: dag.edges,
    given: dag.given,
  });
  fetch(
    'dag',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: data,
    },
  ).then((response) => response.json()).then(
    (questionnaire) => {
      // Display the questionnaire returned by the server
      dag.questionnaire = questionnaire;
      displayQuestionnaire(questionnaire);
    },
  ).catch((error) => {
    console.error('Error:', error);
  });
}
