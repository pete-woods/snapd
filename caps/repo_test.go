// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package caps

import (
	"testing"

	. "gopkg.in/check.v1"

	"github.com/ubuntu-core/snappy/testutil"
)

func TestRepository(t *testing.T) {
	TestingT(t)
}

type RepositorySuite struct {
	// Typical repository with known types
	repo *Repository
	// Empty repository
	emptyRepo *Repository
}

var _ = Suite(&RepositorySuite{})

func (s *RepositorySuite) SetUpTest(c *C) {
	s.repo = NewRepository()
	s.emptyRepo = NewRepository()
	err := LoadBuiltInTypes(s.repo)
	c.Assert(err, IsNil)
}

func (s *RepositorySuite) TestAdd(c *C) {
	cap := &Capability{Name: "name", Label: "label", Type: FileType}
	c.Assert(s.repo.Names(), Not(testutil.Contains), cap.Name)
	err := s.repo.Add(cap)
	c.Assert(err, IsNil)
	c.Assert(s.repo.Names(), DeepEquals, []string{"name"})
	c.Assert(s.repo.Names(), testutil.Contains, cap.Name)
}

func (s *RepositorySuite) TestAddClash(c *C) {
	cap1 := &Capability{Name: "name", Label: "label 1", Type: FileType}
	err := s.repo.Add(cap1)
	c.Assert(err, IsNil)
	cap2 := &Capability{Name: "name", Label: "label 2", Type: FileType}
	err = s.repo.Add(cap2)
	c.Assert(err, ErrorMatches,
		`cannot add capability "name": name already exists`)
	c.Assert(s.repo.Names(), DeepEquals, []string{"name"})
	c.Assert(s.repo.Names(), testutil.Contains, cap1.Name)
}

func (s *RepositorySuite) TestAddInvalidName(c *C) {
	cap := &Capability{Name: "bad-name-", Label: "label", Type: FileType}
	err := s.repo.Add(cap)
	c.Assert(err, ErrorMatches, `"bad-name-" is not a valid snap name`)
	c.Assert(s.repo.Names(), DeepEquals, []string{})
	c.Assert(s.repo.Names(), Not(testutil.Contains), cap.Name)
}

func (s *RepositorySuite) TestAddType(c *C) {
	t := Type("foo")
	err := s.emptyRepo.AddType(t)
	c.Assert(err, IsNil)
	c.Assert(s.emptyRepo.TypeNames(), DeepEquals, []string{"foo"})
	c.Assert(s.emptyRepo.TypeNames(), testutil.Contains, "foo")
}

func (s *RepositorySuite) TestAddTypeClash(c *C) {
	t1 := Type("foo")
	t2 := Type("foo")
	err := s.emptyRepo.AddType(t1)
	c.Assert(err, IsNil)
	err = s.emptyRepo.AddType(t2)
	c.Assert(err, ErrorMatches,
		`cannot add type "foo": name already exists`)
	c.Assert(s.emptyRepo.TypeNames(), DeepEquals, []string{"foo"})
	c.Assert(s.emptyRepo.TypeNames(), testutil.Contains, "foo")
}

func (s *RepositorySuite) TestAddTypeInvalidName(c *C) {
	t := Type("bad-name-")
	err := s.emptyRepo.AddType(t)
	c.Assert(err, ErrorMatches, `"bad-name-" is not a valid snap name`)
	c.Assert(s.emptyRepo.TypeNames(), DeepEquals, []string{})
	c.Assert(s.emptyRepo.TypeNames(), Not(testutil.Contains), string(t))
}

func (s *RepositorySuite) TestRemoveGood(c *C) {
	cap := &Capability{Name: "name", Label: "label", Type: FileType}
	err := s.repo.Add(cap)
	c.Assert(err, IsNil)
	err = s.repo.Remove(cap.Name)
	c.Assert(err, IsNil)
	c.Assert(s.repo.Names(), HasLen, 0)
	c.Assert(s.repo.Names(), Not(testutil.Contains), cap.Name)
}

func (s *RepositorySuite) TestRemoveNoSuchCapability(c *C) {
	err := s.emptyRepo.Remove("name")
	c.Assert(err, ErrorMatches, `can't remove capability "name", no such capability`)
}

func (s *RepositorySuite) TestNames(c *C) {
	// Note added in non-sorted order
	err := s.repo.Add(&Capability{Name: "a", Label: "label-a", Type: FileType})
	c.Assert(err, IsNil)
	err = s.repo.Add(&Capability{Name: "c", Label: "label-c", Type: FileType})
	c.Assert(err, IsNil)
	err = s.repo.Add(&Capability{Name: "b", Label: "label-b", Type: FileType})
	c.Assert(err, IsNil)
	c.Assert(s.repo.Names(), DeepEquals, []string{"a", "b", "c"})
}

func (s *RepositorySuite) TestTypeNames(c *C) {
	c.Assert(s.emptyRepo.TypeNames(), DeepEquals, []string{})
	s.emptyRepo.AddType(Type("a"))
	s.emptyRepo.AddType(Type("b"))
	s.emptyRepo.AddType(Type("c"))
	c.Assert(s.emptyRepo.TypeNames(), DeepEquals, []string{"a", "b", "c"})
}

func (s *RepositorySuite) TestAll(c *C) {
	// Note added in non-sorted order
	err := s.repo.Add(&Capability{Name: "a", Label: "label-a", Type: FileType})
	c.Assert(err, IsNil)
	err = s.repo.Add(&Capability{Name: "c", Label: "label-c", Type: FileType})
	c.Assert(err, IsNil)
	err = s.repo.Add(&Capability{Name: "b", Label: "label-b", Type: FileType})
	c.Assert(err, IsNil)
	c.Assert(s.repo.All(), DeepEquals, []Capability{
		Capability{Name: "a", Label: "label-a", Type: FileType},
		Capability{Name: "b", Label: "label-b", Type: FileType},
		Capability{Name: "c", Label: "label-c", Type: FileType},
	})
}