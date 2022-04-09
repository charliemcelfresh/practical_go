gen:
    # Generate twirp code
	protoc --go_out=. --twirp_out=. rpc/practical_go/service.proto
