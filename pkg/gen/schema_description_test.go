package gen

import (
	"testing"

	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
)

var dummySpan = []int32{1, 0, 1}

func buildDescriptorWithComments(t *testing.T) *descriptorpb.FileDescriptorProto {
	t.Helper()
	return &descriptorpb.FileDescriptorProto{
		Name:    sp("test_comments.proto"),
		Package: sp("testcomments"),
		Syntax:  sp("proto3"),
		EnumType: []*descriptorpb.EnumDescriptorProto{
			{
				Name: sp("Priority"),
				Value: []*descriptorpb.EnumValueDescriptorProto{
					{Name: sp("LOW"), Number: i32p(0)},
					{Name: sp("HIGH"), Number: i32p(1)},
				},
			},
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				// message index 0: MyMessage
				Name: sp("MyMessage"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("title"), Number: i32p(1), Type: ftp(descriptorpb.FieldDescriptorProto_TYPE_STRING), Label: flp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL), JsonName: sp("title")},
					{Name: sp("count"), Number: i32p(2), Type: ftp(descriptorpb.FieldDescriptorProto_TYPE_INT32), Label: flp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL), JsonName: sp("count")},
					{Name: sp("tags"), Number: i32p(3), Type: ftp(descriptorpb.FieldDescriptorProto_TYPE_STRING), Label: flp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED), JsonName: sp("tags")},
					{Name: sp("no_comment"), Number: i32p(4), Type: ftp(descriptorpb.FieldDescriptorProto_TYPE_BOOL), Label: flp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL), JsonName: sp("noComment")},
					{Name: sp("priority"), Number: i32p(5), Type: ftp(descriptorpb.FieldDescriptorProto_TYPE_ENUM), TypeName: sp(".testcomments.Priority"), Label: flp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL), JsonName: sp("priority")},
					{Name: sp("nested"), Number: i32p(6), Type: ftp(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE), TypeName: sp(".testcomments.Inner"), Label: flp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL), JsonName: sp("nested")},
				},
			},
			{
				// message index 1: Inner
				Name: sp("Inner"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("value"), Number: i32p(1), Type: ftp(descriptorpb.FieldDescriptorProto_TYPE_STRING), Label: flp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL), JsonName: sp("value")},
				},
			},
		},
		SourceCodeInfo: &descriptorpb.SourceCodeInfo{
			Location: []*descriptorpb.SourceCodeInfo_Location{
				// Message: path [4, 0] = message_type[0] (MyMessage)
				{Path: []int32{4, 0}, Span: dummySpan, LeadingComments: sp("A test message with various field types")},
				// Fields in MyMessage: path [4, 0, 2, fieldIndex]
				{Path: []int32{4, 0, 2, 0}, Span: dummySpan, LeadingComments: sp("The title of the item")},
				{Path: []int32{4, 0, 2, 1}, Span: dummySpan, LeadingComments: sp("How many items")},
				{Path: []int32{4, 0, 2, 2}, Span: dummySpan, LeadingComments: sp("Tags for categorization")},
				// field index 3 (no_comment) has no entry — tests the no-comment case
				{Path: []int32{4, 0, 2, 4}, Span: dummySpan, LeadingComments: sp("The priority of the task")},
				{Path: []int32{4, 0, 2, 5}, Span: dummySpan, LeadingComments: sp("Nested inner message")},
				// Message: path [4, 1] = message_type[1] (Inner)
				{Path: []int32{4, 1}, Span: dummySpan, LeadingComments: sp("An inner message")},
				{Path: []int32{4, 1, 2, 0}, Span: dummySpan, LeadingComments: sp("The inner value")},
			},
		},
	}
}

func TestDescriptionForDescriptor_FieldComments(t *testing.T) {
	g := NewWithT(t)

	fdp := buildDescriptorWithComments(t)
	file, err := protodesc.NewFile(fdp, nil)
	g.Expect(err).ToNot(HaveOccurred())

	md := file.Messages().Get(0) // MyMessage

	t.Run("scalar field with comment", func(t *testing.T) {
		g := NewWithT(t)
		schema := FieldSchema(md.Fields().ByName("title"), SchemaOptions{})
		g.Expect(schema["description"]).To(Equal("The title of the item"))
	})

	t.Run("integer field with comment", func(t *testing.T) {
		g := NewWithT(t)
		schema := FieldSchema(md.Fields().ByName("count"), SchemaOptions{})
		g.Expect(schema["description"]).To(Equal("How many items"))
	})

	t.Run("repeated field has description on array", func(t *testing.T) {
		g := NewWithT(t)
		schema := FieldSchema(md.Fields().ByName("tags"), SchemaOptions{})
		g.Expect(schema["type"]).To(Equal("array"))
		g.Expect(schema["description"]).To(Equal("Tags for categorization"))
		g.Expect(schema["items"]).ToNot(HaveKey("description"))
	})

	t.Run("field without comment has no description", func(t *testing.T) {
		g := NewWithT(t)
		schema := FieldSchema(md.Fields().ByName("no_comment"), SchemaOptions{})
		g.Expect(schema).ToNot(HaveKey("description"))
	})

	t.Run("enum field gets only field comment", func(t *testing.T) {
		g := NewWithT(t)
		schema := FieldSchema(md.Fields().ByName("priority"), SchemaOptions{})
		g.Expect(schema["description"]).To(Equal("The priority of the task"))
	})

	t.Run("nested message field gets field comment", func(t *testing.T) {
		g := NewWithT(t)
		schema := FieldSchema(md.Fields().ByName("nested"), SchemaOptions{})
		g.Expect(schema["type"]).To(Equal("object"))
		// The nested message schema gets the message-level description,
		// but fieldSchema prepends the field-level comment.
		desc := schema["description"].(string)
		g.Expect(desc).To(ContainSubstring("Nested inner message"))
	})
}

func TestDescriptionForDescriptor_MessageComment(t *testing.T) {
	g := NewWithT(t)

	fdp := buildDescriptorWithComments(t)
	file, err := protodesc.NewFile(fdp, nil)
	g.Expect(err).ToNot(HaveOccurred())

	md := file.Messages().Get(0) // MyMessage
	schema := MessageSchema(md, SchemaOptions{})
	g.Expect(schema["description"]).To(Equal("A test message with various field types"))
}

func TestDescriptionForDescriptor_NoSourceCodeInfo(t *testing.T) {
	g := NewWithT(t)

	// Build a descriptor WITHOUT SourceCodeInfo
	fdp := &descriptorpb.FileDescriptorProto{
		Name:    sp("test_no_sci.proto"),
		Package: sp("testnosci"),
		Syntax:  sp("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: sp("Simple"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{Name: sp("name"), Number: i32p(1), Type: ftp(descriptorpb.FieldDescriptorProto_TYPE_STRING), Label: flp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL), JsonName: sp("name")},
				},
			},
		},
	}
	file, err := protodesc.NewFile(fdp, nil)
	g.Expect(err).ToNot(HaveOccurred())

	md := file.Messages().Get(0)

	schema := MessageSchema(md, SchemaOptions{})
	g.Expect(schema).ToNot(HaveKey("description"))

	fieldSchema := FieldSchema(md.Fields().ByName("name"), SchemaOptions{})
	g.Expect(fieldSchema).ToNot(HaveKey("description"))
}

func TestDescriptionForDescriptor_OpenAIMapMerge(t *testing.T) {
	g := NewWithT(t)

	fdp := &descriptorpb.FileDescriptorProto{
		Name:    sp("test_map_comment.proto"),
		Package: sp("testmapcomment"),
		Syntax:  sp("proto3"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: sp("WithMap"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     sp("labels"),
						Number:   i32p(1),
						Type:     ftp(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE),
						TypeName: sp(".testmapcomment.WithMap.LabelsEntry"),
						Label:    flp(descriptorpb.FieldDescriptorProto_LABEL_REPEATED),
						JsonName: sp("labels"),
					},
				},
				NestedType: []*descriptorpb.DescriptorProto{
					{
						Name: sp("LabelsEntry"),
						Field: []*descriptorpb.FieldDescriptorProto{
							{Name: sp("key"), Number: i32p(1), Type: ftp(descriptorpb.FieldDescriptorProto_TYPE_STRING), Label: flp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL), JsonName: sp("key")},
							{Name: sp("value"), Number: i32p(2), Type: ftp(descriptorpb.FieldDescriptorProto_TYPE_STRING), Label: flp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL), JsonName: sp("value")},
						},
						Options: &descriptorpb.MessageOptions{
							MapEntry: bp(true),
						},
					},
				},
			},
		},
		SourceCodeInfo: &descriptorpb.SourceCodeInfo{
			Location: []*descriptorpb.SourceCodeInfo_Location{
				{Path: []int32{4, 0, 2, 0}, Span: dummySpan, LeadingComments: sp("User-defined labels")},
			},
		},
	}
	file, err := protodesc.NewFile(fdp, nil)
	g.Expect(err).ToNot(HaveOccurred())

	md := file.Messages().Get(0)

	t.Run("standard mode", func(t *testing.T) {
		g := NewWithT(t)
		schema := FieldSchema(md.Fields().ByName("labels"), SchemaOptions{})
		g.Expect(schema["description"]).To(Equal("User-defined labels"))
	})

}

func bp(b bool) *bool { return &b }
