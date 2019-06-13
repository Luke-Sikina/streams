package streams

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strconv"
	"testing"
)

// Steps:
// 1. There exists a file of all the numbers between 0 and 10,000,000
// 2. Stream the file in 1 line at a time
// 3. Filter the lines to something more sane
// 4. Map the lines
// 5. Collect the lines

// Purpose:
// This test is designed to demonstrate that this library can work
// With relatively large data sets. I'm trying to strike a balance between
// scale and reasonable runtime. Hopefully 10 million lines and 1-3 seconds
// represents a good middle ground.
// The test also lays out an example of how I expect this library to be used,
// filtering, mapping and finally reducing / collecting / consuming the stream.
func createStream() Stream {
	file, _ := os.Open("integration_test.txt")
	ch := make(chan interface{}, 1024)

	reader := bufio.NewScanner(file)
	go func() {
		defer close(ch)
		defer closeFile(file)
		for reader.Scan() {
			line := reader.Text()
			num, err := strconv.Atoi(line)
			if err == nil {
				ch <- num
			}
		}
	}()
	return ch
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Printf("Error closing file: %v", err)
	}
}

func DivisibleByTwo(subject interface{}) bool {
	return subject.(int)%2 == 0
}

func DivisibleByThree(subject interface{}) bool {
	return subject.(int)%3 == 0
}

func DivisibleByFive(subject interface{}) bool {
	return subject.(int)%5 == 0
}

func DivisibleBySeven(subject interface{}) bool {
	return subject.(int)%7 == 0
}

func DivisibleByEleven(subject interface{}) bool {
	return subject.(int)%11 == 0
}

func DivisibleByThirteen(subject interface{}) bool {
	return subject.(int)%13 == 0
}

func MapToString(subject interface{}) interface{} {
	return fmt.Sprintf("%v, ", subject)
}

func ReduceConcat(first, second interface{}) interface{} {
	return fmt.Sprintf("%v%v", first, second)
}

func TestStreams(t *testing.T) {
	stream := createStream()
	streams := FromStream(stream, 1024)

	actual := streams.
		Filter(DivisibleByThirteen).
		Filter(DivisibleByEleven).
		Filter(DivisibleBySeven).
		Filter(DivisibleByFive).
		Filter(DivisibleByThree).
		Filter(DivisibleByTwo).
		Map(MapToString).
		Reduce("", ReduceConcat)
	expected := "0, 30030, 60060, 90090, 120120, 150150, 180180, 210210, 240240, 270270, 300300, 330330, 360360, 390390, 420420, 450450, 480480, 510510, 540540, 570570, 600600, 630630, 660660, 690690, 720720, 750750, 780780, 810810, 840840, 870870, 900900, 930930, 960960, 990990, 1021020, 1051050, 1081080, 1111110, 1141140, 1171170, 1201200, 1231230, 1261260, 1291290, 1321320, 1351350, 1381380, 1411410, 1441440, 1471470, 1501500, 1531530, 1561560, 1591590, 1621620, 1651650, 1681680, 1711710, 1741740, 1771770, 1801800, 1831830, 1861860, 1891890, 1921920, 1951950, 1981980, 2012010, 2042040, 2072070, 2102100, 2132130, 2162160, 2192190, 2222220, 2252250, 2282280, 2312310, 2342340, 2372370, 2402400, 2432430, 2462460, 2492490, 2522520, 2552550, 2582580, 2612610, 2642640, 2672670, 2702700, 2732730, 2762760, 2792790, 2822820, 2852850, 2882880, 2912910, 2942940, 2972970, 3003000, 3033030, 3063060, 3093090, 3123120, 3153150, 3183180, 3213210, 3243240, 3273270, 3303300, 3333330, 3363360, 3393390, 3423420, 3453450, 3483480, 3513510, 3543540, 3573570, 3603600, 3633630, 3663660, 3693690, 3723720, 3753750, 3783780, 3813810, 3843840, 3873870, 3903900, 3933930, 3963960, 3993990, 4024020, 4054050, 4084080, 4114110, 4144140, 4174170, 4204200, 4234230, 4264260, 4294290, 4324320, 4354350, 4384380, 4414410, 4444440, 4474470, 4504500, 4534530, 4564560, 4594590, 4624620, 4654650, 4684680, 4714710, 4744740, 4774770, 4804800, 4834830, 4864860, 4894890, 4924920, 4954950, 4984980, 5015010, 5045040, 5075070, 5105100, 5135130, 5165160, 5195190, 5225220, 5255250, 5285280, 5315310, 5345340, 5375370, 5405400, 5435430, 5465460, 5495490, 5525520, 5555550, 5585580, 5615610, 5645640, 5675670, 5705700, 5735730, 5765760, 5795790, 5825820, 5855850, 5885880, 5915910, 5945940, 5975970, 6006000, 6036030, 6066060, 6096090, 6126120, 6156150, 6186180, 6216210, 6246240, 6276270, 6306300, 6336330, 6366360, 6396390, 6426420, 6456450, 6486480, 6516510, 6546540, 6576570, 6606600, 6636630, 6666660, 6696690, 6726720, 6756750, 6786780, 6816810, 6846840, 6876870, 6906900, 6936930, 6966960, 6996990, 7027020, 7057050, 7087080, 7117110, 7147140, 7177170, 7207200, 7237230, 7267260, 7297290, 7327320, 7357350, 7387380, 7417410, 7447440, 7477470, 7507500, 7537530, 7567560, 7597590, 7627620, 7657650, 7687680, 7717710, 7747740, 7777770, 7807800, 7837830, 7867860, 7897890, 7927920, 7957950, 7987980, 8018010, 8048040, 8078070, 8108100, 8138130, 8168160, 8198190, 8228220, 8258250, 8288280, 8318310, 8348340, 8378370, 8408400, 8438430, 8468460, 8498490, 8528520, 8558550, 8588580, 8618610, 8648640, 8678670, 8708700, 8738730, 8768760, 8798790, 8828820, 8858850, 8888880, 8918910, 8948940, 8978970, 9009000, 9039030, 9069060, 9099090, 9129120, 9159150, 9189180, 9219210, 9249240, 9279270, 9309300, 9339330, 9369360, 9399390, 9429420, 9459450, 9489480, 9519510, 9549540, 9579570, 9609600, 9639630, 9669660, 9699690, 9729720, 9759750, 9789780, 9819810, 9849840, 9879870, 9909900, 9939930, 9969960, 9999990, "
	assert.Equal(t, expected, actual)
}
