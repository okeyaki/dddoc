package graph

import (
	"fmt"

	"github.com/awalterschulze/gographviz"
	"github.com/goccy/go-graphviz"
	"github.com/okeyaki/dddoc/lib/platform"
	"github.com/okeyaki/dddoc/lib/platform/source/components"
)

func Render(cs []components.Component) error {
	g, err := newGraph()
	if err != nil {
		return err
	}

	if err := g.Attrs.Add("nodesep", "2.0"); err != nil {
		return err
	}

	if err := renderNodes(g, cs); err != nil {
		return err
	}

	if err := renderEdges(g, cs); err != nil {
		return err
	}

	cg, err := graphviz.ParseBytes([]byte(g.String()))
	if err != nil {
		return err
	}

	return graphviz.New().RenderFilename(
		cg,
		graphviz.PNG,
		platform.GetConfig().GetString("output.file"),
	)
}

func newGraph() (*gographviz.Graph, error) {
	g := gographviz.NewGraph()

	if err := g.SetName("root"); err != nil {
		return nil, err
	}

	if err := g.SetDir(true); err != nil {
		return nil, err
	}

	return g, nil
}

func renderNodes(g *gographviz.Graph, cs []components.Component) error {
	for _, c := range cs {
		options := map[string]string{}
		if err := g.AddNode("root", components.GetComponentName(c), options); err != nil {
			return err
		}
	}

	return nil
}

func renderEdges(g *gographviz.Graph, cs []components.Component) error {
	for _, c := range cs {
		switch c := c.(type) {
		case *components.Entity:
			for _, f := range c.Fields {
				if f.Association == nil {
					continue
				}

				options := map[string]string{
					"label": fmt.Sprintf(`"%s"`, f.Association.Description),
				}
				if err := g.AddEdge(components.GetComponentName(c), f.Association.With, true, options); err != nil {
					return err
				}
			}

		case *components.Factory:
			options := map[string]string{
				"label": fmt.Sprintf(`"%s"`, platform.GetConfig().GetString("output.graph.edge.factory.label")),
				"style": platform.GetConfig().GetString("output.graph.edge.factory.style"),
			}
			if err := g.AddEdge(components.GetComponentName(c), c.Object, true, options); err != nil {
				return err
			}

		case *components.Repository:
			options := map[string]string{
				"label": fmt.Sprintf(`"%s"`, platform.GetConfig().GetString("output.graph.edge.repository.label")),
				"style": platform.GetConfig().GetString("output.graph.edge.repository.style"),
			}
			if err := g.AddEdge(components.GetComponentName(c), c.Object, true, options); err != nil {
				return err
			}

		default:
		}

	}

	return nil
}
