//go:build e2e

package e2e

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// TestSubscription tests the subscription mechanism of the parser.
// The parser processing from a specific block number and waits for the transactions to be observed.
// The transactions `tx0`, `tx1`, `tx2` and `tx3` are present in confirmed ethereum blocks so
// the parser should observe them during the test run.
// The test also ensures that the parser is still running and processing new blocks in parallel.
func (s *ParserTestSuite) TestSubscription() {
	timeout := time.Second * 20

	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	// The following transactions are present in confirmed ethereum blocks:
	tx0 := "0x245295d8C90Fe127932C6fE78daE6D5a4B975098" // present in block 19698125
	tx1 := "0xfb5C635BCC10f3d97a581f11Ba1bdF30F22972C5" // present in block 19698125
	tx2 := "0xf70da97812CB96acDF810712Aa562db8dfA3dbEF" // present in block 19698126
	tx3 := "0xA06769b6d0F7A2BbDe58A81fc888ABCf352C94e6" // present in block 19698130

	s.parser.Subscribe(tx0)
	s.parser.Subscribe(tx1)
	s.parser.Subscribe(tx2)
	s.parser.Subscribe(tx3)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := s.parser.Run(ctx)
		s.Require().NoError(err)
		wg.Done()
	}()

	s.Require().Eventually(func() bool {
		observed0 := s.parser.GetTransactions(tx0)
		observed1 := s.parser.GetTransactions(tx1)
		observed2 := s.parser.GetTransactions(tx2)
		observed3 := s.parser.GetTransactions(tx3)

		allObserved := len(observed0) > 0 && len(observed1) > 0 && len(observed2) > 0 && len(observed3) > 0
		if !allObserved {
			s.T().Logf("observed0: %d, observed1: %d, observed2: %d, observed3: %d",
				len(observed0), len(observed1), len(observed2), len(observed3))
		}

		return allObserved
	}, timeout, time.Millisecond*100)

	s.T().Log("all transactions observed")

	// Ensure that the parser is still running and processing new blocks in parallel.
	s.Require().Eventually(func() bool {
		return s.parser.GetCurrentBlock() > int(s.lastParsedBlock+50)
	}, timeout, time.Millisecond*100)

	cancel()
	wg.Wait()
}

func TestParserTestSuiteSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
