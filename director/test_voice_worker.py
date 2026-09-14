import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from director.voice_worker import VoiceWorker

FAKE = '''import json, os, sys, time
for line in sys.stdin:
 r=json.loads(line)
 if r['text']=='stall': time.sleep(5)
 if r['text']=='crash': sys.exit(1)
 if r['text']=='wrong': r['id']=-1
 print(json.dumps({'id':r['id'],'ok':True}),flush=True)
'''


class VoiceWorkerTests(unittest.TestCase):
    def setUp(self):
        self.folder = tempfile.TemporaryDirectory()
        script = Path(self.folder.name) / 'worker.py'
        script.write_text(FAKE)
        self.worker = VoiceWorker([sys.executable, str(script)], timeout=0.5)

    def tearDown(self):
        self.worker.close()
        self.folder.cleanup()

    def render(self, text):
        self.worker.render(text, None, Path(self.folder.name) / 'voice.wav')

    def test_reuses_live_process_for_different_articles(self):
        self.render('first article')
        process = self.worker.process
        self.render('second article')
        self.assertIs(self.worker.process, process)
        self.assertIsNone(process.poll())

    def test_timeout_terminates_before_next_request(self):
        with self.assertRaises(subprocess.TimeoutExpired):
            self.render('stall')
        self.assertIsNone(self.worker.process)
        self.render('recovered')

    def test_exit_recovers_on_next_request(self):
        with self.assertRaises(OSError):
            self.render('crash')
        self.render('recovered')

    def test_rejects_wrong_response(self):
        with self.assertRaises(OSError):
            self.render('wrong')
        self.assertIsNone(self.worker.process)


if __name__ == '__main__':
    unittest.main()
