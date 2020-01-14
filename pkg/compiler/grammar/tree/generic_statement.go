package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type GenericStatement struct {
	Region input.Region
	Parent Node
	Name *Identifier
	Constraints []TypeName
}

func (statement *GenericStatement) Accept(visitor Visitor) {}

func (statement *GenericStatement) AcceptRecursive(visitor Visitor) {
	statement.Accept(visitor)
	statement.Name.AcceptRecursive(visitor)
	for _, constraint := range statement.Constraints {
		constraint.AcceptRecursive(visitor)
	}
}

func (statement *GenericStatement) Locate() input.Region {
	return statement.Region
}

func (statement *GenericStatement) EnclosingNode() (Node, bool) {
	return statement.Parent, statement.Parent != nil
}

func (statement *GenericStatement) SetEnclosingNode(node Node) {
	statement.Parent = node
}

func (statement *GenericStatement) Matches(target Node) bool {
	if _, isWildcard := target.(*WildcardNode); isWildcard {
		return true
	}
	targetStatement, ok := target.(*GenericStatement)
	return ok && statement.matchesStatement(targetStatement)
}

func (statement *GenericStatement) matchesStatement(target *GenericStatement) bool {
	return statement.Name.Matches(target.Name) &&
		statement.matchesConstraints(target.Constraints)
}

func (statement *GenericStatement) matchesConstraintsOrdered(
	constraints []TypeName) bool {

	if len(statement.Constraints) != len(constraints) {
		return false
	}
	for index, constraint := range statement.Constraints {
		targetConstraint := constraints[index]
		if !constraint.Matches(targetConstraint) {
			return false
		}
	}
	return true
}

func (statement *GenericStatement) matchesConstraints(constraints []TypeName) bool {
	if len(statement.Constraints) != len(constraints) {
		return false
	}
	for _, constraint := range statement.Constraints {
		for _, targetConstraint := range constraints {
			if constraint.Matches(targetConstraint) {
				continue
			}
		}
		return false
	}
	return true
}
